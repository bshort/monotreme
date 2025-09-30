import { Button, Card, Checkbox, Input, Textarea, Typography } from "@mui/joy";
import classnames from "classnames";
import { uniq } from "lodash-es";
import { useEffect, useState, useCallback } from "react";
import { toast } from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router-dom";
import Icon from "@/components/Icon";
import IconUpload from "@/components/IconUpload";
import useLoading from "@/hooks/useLoading";
import useNavigateTo from "@/hooks/useNavigateTo";
import { useShortcutStore, useWorkspaceStore, useUserStore } from "@/stores";
import { getShortcutUpdateMask } from "@/stores/shortcut";
import { Visibility } from "@/types/proto/api/v1/common";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";
import { Role } from "@/types/proto/api/v1/user_service";
import { fetchPageTitle, debounce, generateUrlFriendlyName } from "@/utils/urlMetadata";

interface State {
  shortcutEdit: Shortcut;
}

const ShortcutEditPage = () => {
  const { t } = useTranslation();
  const params = useParams();
  const shortcutUuid = params["shortcutUuid"] || "";
  const navigateTo = useNavigateTo();
  const [state, setState] = useState<State>({
    shortcutEdit: Shortcut.fromPartial({
      visibility: Visibility.WORKSPACE,
      ogMetadata: {
        title: "",
        description: "",
        image: "",
      },
      customIcon: "",
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
  const loadingState = useLoading(true);
  const requestState = useLoading(false);
  const [isFetchingTitle, setIsFetchingTitle] = useState<boolean>(false);
  const [titleWasManuallyEdited, setTitleWasManuallyEdited] = useState<boolean>(false);
  const [nameWasManuallyEdited, setNameWasManuallyEdited] = useState<boolean>(false);

  const shortcut = shortcutStore.getShortcutByUuid(shortcutUuid);
  const havePermission = currentUser.role === Role.ADMIN || shortcut.creatorId === currentUser.id;

  const setPartialState = (partialState: Partial<State>) => {
    setState({
      ...state,
      ...partialState,
    });
  };

  // Helper function to auto-generate name from title
  const autoGenerateNameFromTitle = useCallback((title: string) => {
    const autoGenerateName = currentUser.autoGenerateName ?? true;
    if (nameWasManuallyEdited || !title || !autoGenerateName) return;

    const generatedName = generateUrlFriendlyName(title);
    if (generatedName) {
      setPartialState({
        shortcutEdit: Object.assign(state.shortcutEdit, {
          name: generatedName,
        }),
      });
    }
  }, [nameWasManuallyEdited, state.shortcutEdit, currentUser.autoGenerateName]);

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
            shortcutEdit: Object.assign(state.shortcutEdit, {
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
    }, 500),
    [state.shortcutEdit, titleWasManuallyEdited, autoGenerateNameFromTitle, currentUser.autoGenerateTitle]
  );

  useEffect(() => {
    if (!havePermission) {
      navigateTo("/", { replace: true });
      return;
    }

    const loadShortcut = async () => {
      try {
        const shortcutData = await shortcutStore.getOrFetchShortcutByUuid(shortcutUuid);
        if (shortcutData.id === -1) {
          toast.error("Shortcut not found");
          navigateTo("/", { replace: true });
          return;
        }

        setState({
          shortcutEdit: { ...shortcutData },
        });

        // Set metadata visibility if there's existing OG metadata
        if (shortcutData.ogMetadata?.title || shortcutData.ogMetadata?.description || shortcutData.ogMetadata?.image) {
          setShowOpenGraphMetadata(true);
        }

        loadingState.setFinish();
      } catch (error) {
        console.error("Failed to load shortcut:", error);
        toast.error("Failed to load shortcut");
        navigateTo("/", { replace: true });
      }
    };

    loadShortcut();
  }, [shortcutUuid, havePermission]);

  useEffect(() => {
    workspaceStore.fetchWorkspaceSetting();
  }, []);

  const handleLinkChange = (link: string) => {
    setPartialState({
      shortcutEdit: Object.assign(state.shortcutEdit, {
        link: link,
      }),
    });

    // Trigger debounced title fetch for new URLs
    debouncedFetchTitle(link);
  };

  const handleTitleChange = (title: string) => {
    setTitleWasManuallyEdited(true);
    setPartialState({
      shortcutEdit: Object.assign(state.shortcutEdit, {
        title: title,
      }),
    });

    // Auto-generate name from manually entered title
    autoGenerateNameFromTitle(title);
  };

  const handleNameChange = (name: string) => {
    setNameWasManuallyEdited(true);
    setPartialState({
      shortcutEdit: Object.assign(state.shortcutEdit, {
        name: name,
      }),
    });
  };

  const handleTagValueChange = (tags: string[]) => {
    setPartialState({
      shortcutEdit: Object.assign(state.shortcutEdit, {
        tags: tags,
      }),
    });
  };

  const handleAddTag = () => {
    if (tag.length === 0) return;

    const tags = state.shortcutEdit.tags;
    if (tags.includes(tag)) {
      setTag("");
      return;
    }

    handleTagValueChange([...tags, tag]);
    setTag("");
  };

  const handleRemoveTag = (tagToRemove: string) => {
    handleTagValueChange(state.shortcutEdit.tags.filter((t) => t !== tagToRemove));
  };

  const handleSaveShortcut = async () => {
    if (requestState.isLoading) {
      return;
    }

    const shortcutCreate = state.shortcutEdit;
    if (!shortcutCreate.link) {
      toast.error("Link is required");
      return;
    }

    const originalShortcut = shortcut;
    const updateMask = getShortcutUpdateMask(originalShortcut, shortcutCreate);
    if (updateMask.length === 0) {
      toast.success("Shortcut updated successfully");
      navigateTo(`/shortcut/${shortcutUuid}`);
      return;
    }

    try {
      requestState.setLoading();
      await shortcutStore.updateShortcut(
        {
          id: originalShortcut.id,
          ...shortcutCreate,
        },
        updateMask
      );
      toast.success("Shortcut updated successfully");
      navigateTo(`/shortcut/${shortcutUuid}`);
    } catch (error) {
      console.error(error);
      toast.error("Failed to update shortcut");
    }
    requestState.setFinish();
  };

  const handleCancel = () => {
    navigateTo(`/shortcut/${shortcutUuid}`);
  };

  if (loadingState.isLoading) {
    return (
      <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 pt-4 pb-6 flex flex-col justify-center items-center min-h-screen">
        <Icon.Loader className="w-8 h-8 animate-spin text-gray-500" />
        <p className="mt-4 text-gray-500">{t("common.loading")}</p>
      </div>
    );
  }

  if (!havePermission) {
    return null;
  }

  return (
    <div className="mx-auto max-w-4xl w-full px-4 sm:px-6 md:px-12 pt-4 pb-6 flex flex-col justify-start items-start min-h-screen">
      <div className="w-full mb-6">
        <Typography level="h1" className="mb-2">
          Edit Shortcut
        </Typography>
        <Typography level="body-lg" color="neutral">
          Update your shortcut details and settings
        </Typography>
      </div>

      <Card className="w-full p-6">
        <div className="w-full flex flex-col gap-4">
          {/* Link Field */}
          <div className="w-full flex flex-col gap-2">
            <Typography level="title-sm" className="!font-medium">
              {t("common.link")} *
            </Typography>
            <Input
              placeholder="https://example.com"
              value={state.shortcutEdit.link}
              onChange={(e) => handleLinkChange(e.target.value)}
              startDecorator={<Icon.Link className="w-4 h-auto" />}
              endDecorator={
                isFetchingTitle && (
                  <Icon.Loader className="w-4 h-auto animate-spin text-gray-400" />
                )
              }
            />
          </div>

          {/* Title Field */}
          <div className="w-full flex flex-col gap-2">
            <Typography level="title-sm" className="!font-medium">
              Title
            </Typography>
            <Input
              placeholder="Enter title"
              value={state.shortcutEdit.title || ""}
              onChange={(e) => handleTitleChange(e.target.value)}
              startDecorator={<Icon.Hash className="w-4 h-auto" />}
            />
          </div>

          {/* Name Field */}
          <div className="w-full flex flex-col gap-2">
            <Typography level="title-sm" className="!font-medium">
              Name *
            </Typography>
            <Input
              placeholder="shortcut-name"
              value={state.shortcutEdit.name}
              onChange={(e) => handleNameChange(e.target.value)}
              startDecorator={<span className="text-gray-400 text-sm">{window.location.origin}/s/</span>}
            />
            <Typography level="body-sm" color="neutral">
              This will be used as the shortcut URL
            </Typography>
          </div>

          {/* Description Field */}
          <div className="w-full flex flex-col gap-2">
            <Typography level="title-sm" className="!font-medium">
              Description
            </Typography>
            <Textarea
              placeholder="Enter description"
              value={state.shortcutEdit.description || ""}
              onChange={(e) =>
                setPartialState({
                  shortcutEdit: Object.assign(state.shortcutEdit, {
                    description: e.target.value,
                  }),
                })
              }
              minRows={3}
            />
          </div>

          {/* Custom Icon */}
          <div className="w-full flex flex-col gap-2">
            <Typography level="title-sm" className="!font-medium">
              Custom Icon
            </Typography>
            <IconUpload
              customIcon={state.shortcutEdit.customIcon}
              onIconChange={(icon) =>
                setPartialState({
                  shortcutEdit: Object.assign(state.shortcutEdit, {
                    customIcon: icon,
                  }),
                })
              }
            />
          </div>

          {/* Tags */}
          <div className="w-full flex flex-col gap-2">
            <Typography level="title-sm" className="!font-medium">
              Tags
            </Typography>
            <div className="w-full flex flex-row justify-start items-center flex-wrap gap-2 mb-2">
              {state.shortcutEdit.tags.map((tag) => (
                <div
                  key={tag}
                  className="max-w-xs w-auto px-2 py-1 flex flex-row justify-start items-center bg-gray-200 dark:bg-zinc-700 rounded text-sm cursor-pointer hover:line-through"
                  onClick={() => handleRemoveTag(tag)}
                >
                  <Icon.Tag className="w-3 h-auto mr-1 opacity-60" />
                  <span className="truncate">{tag}</span>
                  <Icon.X className="w-3 h-auto ml-1 opacity-60" />
                </div>
              ))}
            </div>
            <div className="w-full flex flex-row justify-start items-center gap-2">
              <Input
                placeholder="Enter tag"
                value={tag}
                onChange={(e) => setTag(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    handleAddTag();
                  }
                }}
                startDecorator={<Icon.Tag className="w-4 h-auto" />}
                className="flex-1"
              />
              <Button size="sm" onClick={handleAddTag} disabled={!tag.trim()}>
                Add
              </Button>
            </div>
            <div className="flex flex-wrap gap-1 mt-1">
              {tagSuggestions.slice(0, 10).map((suggestion) => (
                <button
                  key={suggestion}
                  className="px-2 py-1 text-xs bg-gray-100 hover:bg-gray-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 rounded transition-colors"
                  onClick={() => {
                    if (!state.shortcutEdit.tags.includes(suggestion)) {
                      handleTagValueChange([...state.shortcutEdit.tags, suggestion]);
                    }
                  }}
                >
                  {suggestion}
                </button>
              ))}
            </div>
          </div>

          {/* Visibility */}
          <div className="w-full flex flex-col gap-2">
            <Typography level="title-sm" className="!font-medium">
              Visibility
            </Typography>
            <div className="w-full flex flex-row justify-start items-center gap-4">
              <div className="flex flex-row justify-start items-center">
                <input
                  type="radio"
                  id="workspace-visibility"
                  name="visibility"
                  checked={state.shortcutEdit.visibility === Visibility.WORKSPACE}
                  onChange={() =>
                    setPartialState({
                      shortcutEdit: Object.assign(state.shortcutEdit, {
                        visibility: Visibility.WORKSPACE,
                      }),
                    })
                  }
                  className="mr-2"
                />
                <label htmlFor="workspace-visibility" className="text-sm">
                  {t("shortcut.visibility.workspace.self")}
                </label>
              </div>
              <div className="flex flex-row justify-start items-center">
                <input
                  type="radio"
                  id="public-visibility"
                  name="visibility"
                  checked={state.shortcutEdit.visibility === Visibility.PUBLIC}
                  onChange={() =>
                    setPartialState({
                      shortcutEdit: Object.assign(state.shortcutEdit, {
                        visibility: Visibility.PUBLIC,
                      }),
                    })
                  }
                  className="mr-2"
                />
                <label htmlFor="public-visibility" className="text-sm">
                  {t("shortcut.visibility.public.self")}
                </label>
              </div>
            </div>
          </div>

          {/* OpenGraph Metadata Toggle */}
          <div className="w-full">
            <Checkbox
              checked={showOpenGraphMetadata}
              onChange={(e) => setShowOpenGraphMetadata(e.target.checked)}
              label="Advanced: Custom Open Graph metadata"
            />
          </div>

          {/* OpenGraph Metadata Fields */}
          {showOpenGraphMetadata && (
            <div className="w-full flex flex-col gap-4 p-4 bg-gray-50 dark:bg-zinc-800/50 rounded-lg">
              <Typography level="title-sm" className="!font-medium">
                Open Graph Metadata
              </Typography>
              <div className="w-full flex flex-col gap-2">
                <Typography level="body-sm">OG Title</Typography>
                <Input
                  placeholder="Custom title for social sharing"
                  value={state.shortcutEdit.ogMetadata?.title || ""}
                  onChange={(e) =>
                    setPartialState({
                      shortcutEdit: Object.assign(state.shortcutEdit, {
                        ogMetadata: {
                          ...state.shortcutEdit.ogMetadata,
                          title: e.target.value,
                        },
                      }),
                    })
                  }
                />
              </div>
              <div className="w-full flex flex-col gap-2">
                <Typography level="body-sm">OG Description</Typography>
                <Textarea
                  placeholder="Custom description for social sharing"
                  value={state.shortcutEdit.ogMetadata?.description || ""}
                  onChange={(e) =>
                    setPartialState({
                      shortcutEdit: Object.assign(state.shortcutEdit, {
                        ogMetadata: {
                          ...state.shortcutEdit.ogMetadata,
                          description: e.target.value,
                        },
                      }),
                    })
                  }
                  minRows={2}
                />
              </div>
              <div className="w-full flex flex-col gap-2">
                <Typography level="body-sm">OG Image URL</Typography>
                <Input
                  placeholder="https://example.com/image.jpg"
                  value={state.shortcutEdit.ogMetadata?.image || ""}
                  onChange={(e) =>
                    setPartialState({
                      shortcutEdit: Object.assign(state.shortcutEdit, {
                        ogMetadata: {
                          ...state.shortcutEdit.ogMetadata,
                          image: e.target.value,
                        },
                      }),
                    })
                  }
                />
              </div>
            </div>
          )}
        </div>
      </Card>

      {/* Action Buttons */}
      <div className="w-full flex flex-row justify-end items-center gap-3 mt-6">
        <Button
          variant="outlined"
          color="neutral"
          onClick={handleCancel}
          disabled={requestState.isLoading}
        >
          Cancel
        </Button>
        <Button
          onClick={handleSaveShortcut}
          disabled={!state.shortcutEdit.link || !state.shortcutEdit.name || requestState.isLoading}
          loading={requestState.isLoading}
        >
          Save Changes
        </Button>
      </div>
    </div>
  );
};

export default ShortcutEditPage;