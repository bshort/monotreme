import { Avatar, Tooltip } from "@mui/joy";
import classNames from "classnames";
import copy from "copy-to-clipboard";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import { absolutifyLink } from "@/helpers/utils";
import { useUserStore, useViewStore, useShortcutStore } from "@/stores";
import { getShortcutUrl } from "@/utils/shortcut";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";
import CustomIcon from "./CustomIcon";
import Icon from "./Icon";
import ShortcutActionsDropdown from "./ShortcutActionsDropdown";
import VisibilityIcon from "./VisibilityIcon";

dayjs.extend(relativeTime);

interface Props {
  shortcut: Shortcut;
  onClick?: () => void;
  dragHandleProps?: any;
}

const ShortcutCard = (props: Props) => {
  const { shortcut, onClick, dragHandleProps } = props;
  const { t } = useTranslation();
  const userStore = useUserStore();
  const viewStore = useViewStore();
  const shortcutStore = useShortcutStore();
  const creator = userStore.getUserById(shortcut.creatorId);
  const shortcutLink = absolutifyLink(getShortcutUrl(shortcut.name));
  const [isTogglingFavorite, setIsTogglingFavorite] = useState(false);

  useEffect(() => {
    userStore.getOrFetchUserById(shortcut.creatorId);
  }, []);

  const handleCopyButtonClick = () => {
    copy(shortcutLink);
    toast.success("Shortcut link copied to clipboard.");
  };

  const handleToggleFavorite = async (e: React.MouseEvent) => {
    e.stopPropagation();
    if (isTogglingFavorite) return;

    setIsTogglingFavorite(true);
    try {
      await shortcutStore.updateShortcut(
        {
          id: shortcut.id,
          isFavorite: !shortcut.isFavorite,
        },
        ["is_favorite"]
      );
      toast.success(shortcut.isFavorite ? "Removed from favorites" : "Added to favorites");
    } catch (error) {
      toast.error("Failed to update favorite status");
      console.error("Error toggling favorite:", error);
    } finally {
      setIsTogglingFavorite(false);
    }
  };

  return (
    <div
      className={classNames(
        "group px-4 py-3 w-full flex flex-col justify-start items-start border rounded-lg hover:shadow dark:border-zinc-700",
        !dragHandleProps && "cursor-pointer",
      )}
      onClick={dragHandleProps ? undefined : onClick}
    >
      <div className="w-full flex flex-row justify-between items-center">
        {dragHandleProps && (
          <div
            {...dragHandleProps}
            className="mr-2 cursor-grab active:cursor-grabbing text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 flex items-center"
            title={`Drag to reorder (User Order: ${shortcut.userOrder})`}
          >
            <Icon.GripVertical className="w-5 h-5" />
          </div>
        )}
        <div className={classNames("flex flex-row justify-start items-center mr-1 shrink-0", dragHandleProps ? "w-[calc(100%-36px)]" : "w-[calc(100%-16px)]")}>
          <Link
            className={classNames("w-8 h-8 flex justify-center items-center overflow-clip shrink-0")}
            to={`/shortcut/${shortcut.uuid}`}
            viewTransition
            onClick={(e) => e.stopPropagation()}
          >
            <CustomIcon customIcon={shortcut.customIcon} url={shortcut.link} />
          </Link>
          <div className="ml-2 w-[calc(100%-24px)] flex flex-col justify-start items-start">
            <div className="w-full flex flex-row justify-start items-center leading-tight">
              <Tooltip title={shortcut.isFavorite ? "Remove from favorites" : "Add to favorites"} variant="solid" placement="top" arrow>
                <button
                  className={classNames(
                    "cursor-pointer mr-1 shrink-0 transition-all",
                    shortcut.isFavorite ? "text-yellow-500 hover:text-yellow-600" : "text-gray-300 hover:text-yellow-400",
                    isTogglingFavorite && "opacity-50 cursor-not-allowed"
                  )}
                  onClick={handleToggleFavorite}
                  disabled={isTogglingFavorite}
                >
                  <Icon.Star
                    className={classNames(
                      "w-4 h-auto transition-all",
                      shortcut.isFavorite && "fill-yellow-500 stroke-yellow-600"
                    )}
                    style={shortcut.isFavorite ? { fill: '#eab308', stroke: '#ca8a04' } : {}}
                  />
                </button>
              </Tooltip>
              <a
                className={classNames(
                  "max-w-[calc(100%-60px)] flex flex-row justify-start items-center mr-1 cursor-pointer hover:opacity-80 hover:underline",
                )}
                target="_blank"
                href={shortcutLink}
                onClick={(e) => e.stopPropagation()}
              >
                <div className="truncate">
                  <span className="dark:text-gray-400">{shortcut.title}</span>
                  {shortcut.title ? (
                    <span className="text-gray-500">({shortcut.name})</span>
                  ) : (
                    <span className="truncate dark:text-gray-400">{shortcut.name}</span>
                  )}
                </div>
                <span className="hidden group-hover:block ml-1 cursor-pointer shrink-0">
                  <Icon.ExternalLink className="w-4 h-auto text-gray-600" />
                </span>
              </a>
              <Tooltip title="Copy" variant="solid" placement="top" arrow>
                <button
                  className="hidden group-hover:block cursor-pointer text-gray-500 hover:opacity-80"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleCopyButtonClick();
                  }}
                >
                  <Icon.Clipboard className="w-4 h-auto mx-auto" />
                </button>
              </Tooltip>
            </div>
            <div className="pr-4 leading-tight w-full text-sm">
              <a
                className="block truncate text-gray-400 dark:text-gray-500 hover:underline"
                href={shortcutLink}
                target="_blank"
                onClick={(e) => e.stopPropagation()}
              >
                {shortcut.link}
              </a>
            </div>
          </div>
        </div>
        <div className="h-full pt-2 flex flex-row justify-end items-start">
          <ShortcutActionsDropdown shortcut={shortcut} />
        </div>
      </div>
      <div className="mt-2 w-full flex flex-row justify-start items-start gap-2 truncate">
        {shortcut.tags.map((tag) => {
          return (
            <span
              key={tag}
              className="max-w-[8rem] truncate text-gray-400 dark:text-gray-500 text-sm leading-4 cursor-pointer hover:opacity-80"
              onClick={(e) => {
                e.stopPropagation();
                viewStore.setFilter({ tag: tag });
              }}
            >
              #{tag}
            </span>
          );
        })}
        {shortcut.tags.length === 0 && <span className="text-gray-400 text-sm leading-4 italic">No tags</span>}
      </div>
      <div className="w-full mt-2 flex gap-2 overflow-x-auto">
        <Tooltip title={creator.nickname} variant="solid" placement="top" arrow>
          <Avatar
            className="dark:bg-zinc-800"
            sx={{
              "--Avatar-size": "24px",
            }}
            alt={creator.nickname.toUpperCase()}
          ></Avatar>
        </Tooltip>
        <Tooltip title={t(`shortcut.visibility.${shortcut.visibility.toLowerCase()}.description`)} variant="solid" placement="top" arrow>
          <div
            className="w-auto leading-5 flex flex-row justify-start items-center flex-nowrap whitespace-nowrap cursor-pointer text-gray-400 text-sm"
            onClick={(e) => {
              e.stopPropagation();
              viewStore.setFilter({ visibility: shortcut.visibility });
            }}
          >
            <VisibilityIcon className="w-4 h-auto mr-1 opacity-70" visibility={shortcut.visibility} />
            {t(`shortcut.visibility.${shortcut.visibility.toLowerCase()}.self`)}
          </div>
        </Tooltip>
        <Tooltip title="View count" variant="solid" placement="top" arrow>
          <Link
            className="w-auto leading-5 flex flex-row justify-start items-center flex-nowrap whitespace-nowrap cursor-pointer text-gray-400 text-sm"
            to={`/shortcut/${shortcut.uuid}#analytics`}
            viewTransition
            onClick={(e) => e.stopPropagation()}
          >
            <Icon.BarChart2 className="w-4 h-auto mr-1 opacity-70" />
            {t("shortcut.visits", { count: shortcut.viewCount })}
          </Link>
        </Tooltip>
      </div>
      <div className="w-full mt-2 flex flex-col gap-1 text-xs text-gray-400 dark:text-gray-500">
        {shortcut.createdTime && (
          <Tooltip title={dayjs(shortcut.createdTime).format("YYYY-MM-DD HH:mm:ss")} variant="solid" placement="top" arrow>
            <div className="flex flex-row items-center">
              <Icon.Calendar className="w-3.5 h-auto mr-1.5 opacity-70" />
              <span>Created {dayjs(shortcut.createdTime).fromNow()}</span>
            </div>
          </Tooltip>
        )}
        {shortcut.updatedTime && (
          <Tooltip title={dayjs(shortcut.updatedTime).format("YYYY-MM-DD HH:mm:ss")} variant="solid" placement="top" arrow>
            <div className="flex flex-row items-center">
              <Icon.Clock className="w-3.5 h-auto mr-1.5 opacity-70" />
              <span>Updated {dayjs(shortcut.updatedTime).fromNow()}</span>
            </div>
          </Tooltip>
        )}
      </div>
    </div>
  );
};

export default ShortcutCard;
