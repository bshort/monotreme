import { Avatar, Tooltip } from "@mui/joy";
import classNames from "classnames";
import copy from "copy-to-clipboard";
import { useEffect } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import { absolutifyLink } from "@/helpers/utils";
import { useUserStore, useViewStore } from "@/stores";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";
import Icon from "./Icon";
import LinkFavicon from "./LinkFavicon";
import ShortcutActionsDropdown from "./ShortcutActionsDropdown";
import VisibilityIcon from "./VisibilityIcon";

interface Props {
  shortcut: Shortcut;
  className?: string;
  showActions?: boolean;
  onClick?: () => void;
}

const ShortcutListView = (props: Props) => {
  const { shortcut, className, showActions, onClick } = props;
  const { t } = useTranslation();
  const userStore = useUserStore();
  const viewStore = useViewStore();
  const creator = userStore.getUserById(shortcut.creatorId);
  const shortcutLink = absolutifyLink(`/s/${shortcut.name}`);

  useEffect(() => {
    userStore.getOrFetchUserById(shortcut.creatorId);
  }, []);

  const handleCopyButtonClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    copy(shortcutLink);
    toast.success("Shortcut link copied to clipboard.");
  };

  return (
    <div
      className={classNames(
        "group w-full px-4 py-3 flex flex-col justify-start items-start border rounded-lg hover:bg-gray-50 dark:border-zinc-800 dark:hover:bg-zinc-800/50 cursor-pointer",
        className,
      )}
      onClick={onClick}
    >
      {/* First line: favicon, title, tags */}
      <div className="w-full flex flex-row justify-between items-start">
        <div className="flex flex-row justify-start items-center flex-1 min-w-0">
          <div className="w-5 h-5 flex justify-center items-center overflow-clip shrink-0 mr-3">
            <LinkFavicon url={shortcut.link} />
          </div>
          <div className="flex flex-row justify-start items-center min-w-0 mr-3">
            {shortcut.title ? (
              <>
                <span className="dark:text-gray-300 font-medium mr-2 truncate">{shortcut.title}</span>
                <span className="text-gray-500 text-sm">({shortcut.name})</span>
              </>
            ) : (
              <span className="dark:text-gray-300 font-medium truncate">{shortcut.name}</span>
            )}
          </div>
          {/* Tags */}
          <div className="flex flex-row justify-start items-center gap-1 min-w-0">
            {shortcut.tags.slice(0, 3).map((tag) => (
              <span
                key={tag}
                className="max-w-[6rem] truncate text-blue-600 dark:text-blue-400 text-xs px-1.5 py-0.5 bg-blue-50 dark:bg-blue-900/30 rounded cursor-pointer hover:opacity-80"
                onClick={(e) => {
                  e.stopPropagation();
                  viewStore.setFilter({ tag: tag });
                }}
              >
                #{tag}
              </span>
            ))}
            {shortcut.tags.length > 3 && (
              <span className="text-gray-400 text-xs">+{shortcut.tags.length - 3}</span>
            )}
          </div>
        </div>
        {showActions && (
          <div className="flex flex-row justify-end items-center shrink-0 ml-2" onClick={(e) => e.stopPropagation()}>
            <Tooltip title="Copy" variant="solid" placement="top" arrow>
              <button
                className="hidden group-hover:block cursor-pointer text-gray-500 hover:opacity-80 mr-1"
                onClick={handleCopyButtonClick}
              >
                <Icon.Clipboard className="w-4 h-auto" />
              </button>
            </Tooltip>
            <ShortcutActionsDropdown shortcut={shortcut} />
          </div>
        )}
      </div>

      {/* Second line: full URL and metrics */}
      <div className="w-full mt-2 flex flex-row justify-between items-center">
        <div className="flex flex-row justify-start items-center flex-1 min-w-0 mr-4">
          <a
            className="truncate text-sm text-gray-500 dark:text-gray-400 hover:underline hover:text-gray-700 dark:hover:text-gray-300"
            href={shortcut.link}
            target="_blank"
            onClick={(e) => e.stopPropagation()}
          >
            {shortcut.link}
          </a>
        </div>

        {/* Metrics */}
        <div className="flex flex-row justify-end items-center gap-3 shrink-0">
          <Tooltip title={creator.nickname} variant="solid" placement="top" arrow>
            <Avatar
              className="dark:bg-zinc-700"
              sx={{
                "--Avatar-size": "20px",
                fontSize: "10px",
              }}
              alt={creator.nickname.toUpperCase()}
            />
          </Tooltip>

          <Tooltip title={t(`shortcut.visibility.${shortcut.visibility.toLowerCase()}.description`)} variant="solid" placement="top" arrow>
            <div
              className="flex flex-row justify-start items-center cursor-pointer text-gray-400 text-xs hover:opacity-80"
              onClick={(e) => {
                e.stopPropagation();
                viewStore.setFilter({ visibility: shortcut.visibility });
              }}
            >
              <VisibilityIcon className="w-3 h-auto mr-1 opacity-70" visibility={shortcut.visibility} />
              {t(`shortcut.visibility.${shortcut.visibility.toLowerCase()}.self`)}
            </div>
          </Tooltip>

          <Tooltip title="View count" variant="solid" placement="top" arrow>
            <Link
              className="flex flex-row justify-start items-center cursor-pointer text-gray-400 text-xs hover:opacity-80"
              to={`/shortcut/${shortcut.id}#analytics`}
              viewTransition
              onClick={(e) => e.stopPropagation()}
            >
              <Icon.BarChart2 className="w-3 h-auto mr-1 opacity-70" />
              {t("shortcut.visits", { count: shortcut.viewCount })}
            </Link>
          </Tooltip>
        </div>
      </div>
    </div>
  );
};

export default ShortcutListView;