import { Button, Checkbox, DialogActions, DialogContent, DialogTitle, Divider, Drawer, Input, ModalClose, Textarea } from "@mui/joy";
import classnames from "classnames";
import { isUndefined, uniq } from "lodash-es";
import { useEffect, useState, useCallback } from "react";
import { toast } from "react-hot-toast";
import { useTranslation } from "react-i18next";
import useLoading from "@/hooks/useLoading";
import { useShortcutStore, useWorkspaceStore, useUserStore } from "@/stores";
import { getShortcutUpdateMask } from "@/stores/shortcut";
import { Visibility } from "@/types/proto/api/v1/common";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";
import { fetchPageTitle, debounce, generateUrlFriendlyName } from "@/utils/urlMetadata";
import Icon from "./Icon";
import IconUpload from "./IconUpload";

interface Props {
  shortcutId?: number;
  initialShortcut?: Partial<Shortcut>;
  onClose: () => void;
  onConfirm?: () => void;
}

interface State {
  shortcutCreate: Shortcut;
}

const CreateShortcutDrawer: React.FC<Props> = (props: Props) => {
  const { onClose, onConfirm, shortcutId, initialShortcut } = props;
  const { t } = useTranslation();
  const [state, setState] = useState<State>({
    shortcutCreate: Shortcut.fromPartial({
      visibility: Visibility.WORKSPACE,
      ogMetadata: {
        title: "",
        description: "",
        image: "",
      },
      customIcon: "",
      ...initialShortcut,
    }),
  });
  const shortcutStore = useShortcutStore();
  const workspaceStore = useWorkspaceStore();
  const userStore = useUserStore();
  const currentUser = userStore.getCurrentUser();
  const [showOpenGraphMetadata, setShowOpenGraphMetadata] = useState<boolean>(false);
  const shortcutList = shortcutStore.getShortcutList();
  const [tag, setTag] = useState<string>("");
  const tagSuggestions = uniq(shortcutList.map((shortcut) => shortcut.tags).flat());
  const isCreating = isUndefined(shortcutId);
  const loadingState = useLoading(!isCreating);
  const requestState = useLoading(false);
  const [isFetchingTitle, setIsFetchingTitle] = useState<boolean>(false);
  const [titleWasManuallyEdited, setTitleWasManuallyEdited] = useState<boolean>(false);
  const [nameWasManuallyEdited, setNameWasManuallyEdited] = useState<boolean>(false);

  const setPartialState = (partialState: Partial<State>) => {
    setState({
      ...state,
      ...partialState,
    });
  };

  // Helper function to auto-generate name from title
  const autoGenerateNameFromTitle = useCallback((title: string) => {
    const autoGenerateName = currentUser.autoGenerateName ?? true;
    if (!isCreating || nameWasManuallyEdited || !title || !autoGenerateName) return;

    const generatedName = generateUrlFriendlyName(title);
    if (generatedName) {
      setPartialState({
        shortcutCreate: Object.assign(state.shortcutCreate, {
          name: generatedName,
        }),
      });
    }
  }, [isCreating, nameWasManuallyEdited, state.shortcutCreate, currentUser.autoGenerateName]);

  // Create debounced function for fetching page title
  const debouncedFetchTitle = useCallback(
    debounce(async (url: string) => {
      const autoGenerateTitle = currentUser.autoGenerateTitle ?? true;
      if (!url || titleWasManuallyEdited || !autoGenerateTitle) return;

      setIsFetchingTitle(true);
      try {
        const title = await fetchPageTitle(url);
        if (title && !titleWasManuallyEdited) {
          setPartialState({
            shortcutCreate: Object.assign(state.shortcutCreate, {
              title: title,
            }),
          });
          // Also auto-generate name from the fetched title
          autoGenerateNameFromTitle(title);
        }
      } catch (error) {
        console.warn('Failed to fetch page title:', error);
      } finally {
        setIsFetchingTitle(false);
      }
    }, 1000), // 1 second delay
    [titleWasManuallyEdited, state.shortcutCreate, autoGenerateNameFromTitle, currentUser.autoGenerateTitle]
  );

  useEffect(() => {
    // Set default visibility based on user preference, then workspace setting
    let defaultVisibility = Visibility.WORKSPACE;

    if (currentUser.defaultVisibility) {
      // Map string to enum
      switch (currentUser.defaultVisibility) {
        case "PUBLIC":
          defaultVisibility = Visibility.PUBLIC;
          break;
        case "PRIVATE":
          defaultVisibility = Visibility.PRIVATE;
          break;
        case "WORKSPACE":
        default:
          defaultVisibility = Visibility.WORKSPACE;
          break;
      }
    } else if (workspaceStore.setting.defaultVisibility !== Visibility.VISIBILITY_UNSPECIFIED) {
      defaultVisibility = workspaceStore.setting.defaultVisibility;
    }

    setPartialState({
      shortcutCreate: Object.assign(state.shortcutCreate, {
        visibility: defaultVisibility,
      }),
    });
  }, [currentUser.defaultVisibility]);

  useEffect(() => {
    if (shortcutId) {
      const shortcut = shortcutStore.getShortcutById(shortcutId);
      if (shortcut) {
        setState({
          ...state,
          shortcutCreate: Object.assign(state.shortcutCreate, {
            name: shortcut.name,
            link: shortcut.link,
            title: shortcut.title,
            description: shortcut.description,
            visibility: shortcut.visibility,
            ogMetadata: shortcut.ogMetadata,
            customIcon: shortcut.customIcon,
          }),
        });
        setTag(shortcut.tags.join(" "));
        // When editing an existing shortcut, consider both title and name as manually edited
        // so we don't override them with auto-generated values
        setTitleWasManuallyEdited(true);
        setNameWasManuallyEdited(true);
        loadingState.setFinish();
      }
    }
  }, [shortcutId]);

  if (loadingState.isLoading) {
    return null;
  }

  const handleNameInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    // Mark that the name has been manually edited
    setNameWasManuallyEdited(true);

    setPartialState({
      shortcutCreate: Object.assign(state.shortcutCreate, {
        name: e.target.value.replace(/\s+/g, "-"),
      }),
    });
  };

  const handleLinkInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newLink = e.target.value;
    setPartialState({
      shortcutCreate: Object.assign(state.shortcutCreate, {
        link: newLink,
      }),
    });

    // Only fetch title for new shortcuts and if title hasn't been manually edited and user preference allows it
    const autoGenerateTitle = currentUser.autoGenerateTitle ?? true;
    if (isCreating && !titleWasManuallyEdited && newLink.trim() && autoGenerateTitle) {
      debouncedFetchTitle(newLink.trim());
    }
  };

  const handleTitleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    // Mark that the title has been manually edited
    setTitleWasManuallyEdited(true);

    const newTitle = e.target.value;
    setPartialState({
      shortcutCreate: Object.assign(state.shortcutCreate, {
        title: newTitle,
      }),
    });

    // Auto-generate name from manually entered title
    autoGenerateNameFromTitle(newTitle);
  };

  const handleDescriptionInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPartialState({
      shortcutCreate: Object.assign(state.shortcutCreate, {
        description: e.target.value,
      }),
    });
  };

  const handleTagsInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setTag(text);
  };

  const handleOpenGraphMetadataImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPartialState({
      shortcutCreate: Object.assign(state.shortcutCreate, {
        ogMetadata: {
          ...state.shortcutCreate.ogMetadata,
          image: e.target.value,
        },
      }),
    });
  };

  const handleOpenGraphMetadataTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPartialState({
      shortcutCreate: Object.assign(state.shortcutCreate, {
        ogMetadata: {
          ...state.shortcutCreate.ogMetadata,
          title: e.target.value,
        },
      }),
    });
  };

  const handleOpenGraphMetadataDescriptionChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setPartialState({
      shortcutCreate: Object.assign(state.shortcutCreate, {
        ogMetadata: {
          ...state.shortcutCreate.ogMetadata,
          description: e.target.value,
        },
      }),
    });
  };

  const handleTagSuggestionsClick = (suggestion: string) => {
    if (tag === "") {
      setTag(suggestion);
    } else {
      setTag(`${tag} ${suggestion}`);
    }
  };

  const handleCustomIconChange = (iconData: string) => {
    setPartialState({
      shortcutCreate: Object.assign(state.shortcutCreate, {
        customIcon: iconData,
      }),
    });
  };

  const handleSaveBtnClick = async () => {
    if (!state.shortcutCreate.name || !state.shortcutCreate.link) {
      toast.error("Please fill in required fields.");
      return;
    }

    try {
      const tags = tag.split(" ").filter(Boolean);
      if (shortcutId) {
        const originShortcut = shortcutStore.getShortcutById(shortcutId);
        const updatingShortcut = {
          ...state.shortcutCreate,
          id: shortcutId,
          tags,
        };
        await shortcutStore.updateShortcut(updatingShortcut, getShortcutUpdateMask(originShortcut, updatingShortcut));
      } else {
        await shortcutStore.createShortcut({
          ...state.shortcutCreate,
          tags,
        });
      }

      if (onConfirm) {
        onConfirm();
      } else {
        onClose();
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.details);
    }
  };

  return (
    <Drawer anchor="right" open={true} onClose={onClose}>
      <DialogTitle>{isCreating ? "Create Shortcut" : "Edit Shortcut"}</DialogTitle>
      <ModalClose />
      <DialogContent className="w-full max-w-full">
        <div className="overflow-y-auto w-full mt-2 px-4 pb-4 sm:w-[24rem]">
          <div className="w-full flex flex-col justify-start items-start mb-3">
            <span className="mb-2">
              Link <span className="text-red-600">*</span>
            </span>
            <Input
              className="w-full"
              type="text"
              placeholder="The destination link of the shortcut"
              value={state.shortcutCreate.link}
              onChange={handleLinkInputChange}
            />
          </div>
          <div className="w-full flex flex-col justify-start items-start mb-3">
            <span className="mb-2">
              Name <span className="text-red-600">*</span>
              {isCreating && !nameWasManuallyEdited && (currentUser.autoGenerateName ?? true) && (
                <span className="text-xs text-gray-500 ml-2">(auto-generated from title)</span>
              )}
            </span>
            <Input
              className="w-full"
              type="text"
              startDecorator="s/"
              placeholder={
                isCreating && !nameWasManuallyEdited && (currentUser.autoGenerateName ?? true)
                  ? "Will be generated from title..."
                  : "An easy name to remember"
              }
              value={state.shortcutCreate.name}
              onChange={handleNameInputChange}
            />
          </div>
          <div className="w-full flex flex-col justify-start items-start mb-3">
            <span className="mb-2">Title</span>
            <Input
              className="w-full"
              type="text"
              placeholder={
                isFetchingTitle
                  ? "Fetching page title..."
                  : (currentUser.autoGenerateTitle ?? true) && isCreating
                  ? "Auto-filled from URL or enter manually"
                  : "The title of the shortcut"
              }
              value={state.shortcutCreate.title}
              onChange={handleTitleInputChange}
              endDecorator={
                isFetchingTitle ? (
                  <Icon.Loader className="w-4 h-4 animate-spin text-gray-400" />
                ) : null
              }
            />
          </div>
          <div className="w-full flex flex-col justify-start items-start mb-3">
            <span className="mb-2">Description</span>
            <Input
              className="w-full"
              type="text"
              placeholder="A short description of the shortcut"
              value={state.shortcutCreate.description}
              onChange={handleDescriptionInputChange}
            />
          </div>
          <div className="w-full flex flex-col justify-start items-start mb-3">
            <span className="mb-2">Custom Icon</span>
            <IconUpload
              value={state.shortcutCreate.customIcon}
              onChange={handleCustomIconChange}
              fallbackUrl={state.shortcutCreate.link}
              className="w-16 h-16"
            />
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Upload a custom icon (PNG, JPG, ICO) or leave empty to use favicon
            </p>
          </div>
          <div className="w-full flex flex-col justify-start items-start mb-3">
            <span className="mb-2">Tags</span>
            <Input className="w-full" type="text" placeholder="The tags of shortcut" value={tag} onChange={handleTagsInputChange} />
            {tagSuggestions.length > 0 && (
              <div className="w-full flex flex-row justify-start items-start mt-2">
                <Icon.Asterisk className="w-4 h-auto shrink-0 mx-1 text-gray-400 dark:text-gray-500" />
                <div className="w-auto flex flex-row justify-start items-start flex-wrap gap-x-2 gap-y-1">
                  {tagSuggestions.map((tag) => (
                    <span
                      className="text-gray-600 dark:text-gray-500 cursor-pointer max-w-[6rem] truncate block text-sm flex-nowrap leading-4 hover:text-black dark:hover:text-gray-400"
                      key={tag}
                      onClick={() => handleTagSuggestionsClick(tag)}
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
          <div className="w-full flex flex-col justify-start items-start mb-3">
            <Checkbox
              className="w-full dark:text-gray-400"
              checked={state.shortcutCreate.visibility === Visibility.PUBLIC}
              label={t(`shortcut.visibility.public.description`)}
              onChange={(e) =>
                setPartialState({
                  shortcutCreate: Object.assign(state.shortcutCreate, {
                    visibility: e.target.checked ? Visibility.PUBLIC : Visibility.WORKSPACE,
                  }),
                })
              }
            />
          </div>
          <Divider className="text-gray-500">More</Divider>
          <div className="w-full flex flex-col justify-start items-start border rounded-md mt-3 overflow-hidden dark:border-zinc-800">
            <div
              className={classnames(
                "w-full flex flex-row justify-between items-center px-2 py-1 cursor-pointer hover:bg-gray-100 dark:hover:bg-zinc-800",
                showOpenGraphMetadata ? "bg-gray-100 border-b dark:bg-zinc-800 dark:border-b-zinc-700" : "",
              )}
              onClick={() => setShowOpenGraphMetadata(!showOpenGraphMetadata)}
            >
              <span className="text-sm flex flex-row justify-start items-center">
                Social media metadata
                <Icon.Sparkles className="w-4 h-auto shrink-0 ml-1 text-blue-600 dark:text-blue-500" />
              </span>
              <button className="w-7 h-7 p-1 rounded-md">
                <Icon.ChevronDown className={classnames("w-4 h-auto text-gray-500", showOpenGraphMetadata ? "transform rotate-180" : "")} />
              </button>
            </div>
            {showOpenGraphMetadata && (
              <div className="w-full px-2 py-1">
                <div className="w-full flex flex-col justify-start items-start mb-3">
                  <span className="mb-2 text-sm">Image URL</span>
                  <Input
                    className="w-full"
                    type="text"
                    placeholder="https://the.link.to/the/image.png"
                    size="sm"
                    value={state.shortcutCreate.ogMetadata?.image}
                    onChange={handleOpenGraphMetadataImageChange}
                  />
                </div>
                <div className="w-full flex flex-col justify-start items-start mb-3">
                  <span className="mb-2 text-sm">Title</span>
                  <Input
                    className="w-full"
                    type="text"
                    placeholder="Monotreme - An open source, self-hosted platform for sharing and managing your most frequently used links"
                    size="sm"
                    value={state.shortcutCreate.ogMetadata?.title}
                    onChange={handleOpenGraphMetadataTitleChange}
                  />
                </div>
                <div className="w-full flex flex-col justify-start items-start mb-3">
                  <span className="mb-2 text-sm">Description</span>
                  <Textarea
                    className="w-full"
                    placeholder="An open source, self-hosted platform for sharing and managing your most frequently used links."
                    size="sm"
                    maxRows={3}
                    value={state.shortcutCreate.ogMetadata?.description}
                    onChange={handleOpenGraphMetadataDescriptionChange}
                  />
                </div>
              </div>
            )}
          </div>
        </div>
      </DialogContent>
      <DialogActions>
        <div className="w-full flex flex-row justify-end items-center px-3 py-4 space-x-2">
          <Button color="neutral" variant="plain" disabled={requestState.isLoading} loading={requestState.isLoading} onClick={onClose}>
            {t("common.cancel")}
          </Button>
          <Button color="primary" disabled={requestState.isLoading} loading={requestState.isLoading} onClick={handleSaveBtnClick}>
            {t("common.save")}
          </Button>
        </div>
      </DialogActions>
    </Drawer>
  );
};

export default CreateShortcutDrawer;
