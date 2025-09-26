import classNames from "classnames";
import copy from "copy-to-clipboard";
import toast from "react-hot-toast";
import { Link } from "react-router-dom";
import { absolutifyLink } from "@/helpers/utils";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";
import { getShortcutUrl } from "@/utils/shortcut";
import CustomIcon from "./CustomIcon";
import Icon from "./Icon";
import ShortcutActionsDropdown from "./ShortcutActionsDropdown";

interface Props {
  shortcut: Shortcut;
  className?: string;
  showActions?: boolean;
  alwaysShowLink?: boolean;
  onClick?: () => void;
}

const ShortcutView = (props: Props) => {
  const { shortcut, className, showActions, alwaysShowLink, onClick } = props;
  const shortcutLink = absolutifyLink(getShortcutUrl(shortcut.name));

  const handleCopyButtonClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    copy(shortcutLink);
    toast.success("Shortcut link copied to clipboard.");
  };

  return (
    <div
      className={classNames(
        "group w-full px-3 py-2 flex flex-col justify-start items-start border rounded-lg hover:bg-gray-100 dark:border-zinc-800 dark:hover:bg-zinc-800",
        className,
      )}
      onClick={onClick}
    >
      {/* First row: icon, title, actions */}
      <div className="w-full flex flex-row justify-start items-center">
        <div className={classNames("w-5 h-5 flex justify-center items-center overflow-clip shrink-0")}>
          <CustomIcon customIcon={shortcut.customIcon} url={shortcut.link} />
        </div>
        <div className="ml-2 w-full truncate">
          {shortcut.title ? (
            <>
              <span className="dark:text-gray-400">{shortcut.title}</span>
              <span className="text-gray-500">({shortcut.name})</span>
            </>
          ) : (
            <>
              <span className="dark:text-gray-400">{shortcut.name}</span>
            </>
          )}
        </div>
        <div className="ml-1 flex flex-row justify-end items-center shrink-0">
          <button
            className={classNames(
              "hidden group-hover:block mr-1 w-6 h-6 p-1 shrink-0 rounded-lg bg-gray-200 dark:bg-zinc-900 hover:opacity-80",
            )}
            onClick={handleCopyButtonClick}
            title="Copy shortcut link"
          >
            <Icon.Clipboard className="w-3 h-auto text-gray-400 shrink-0" />
          </button>
          <Link
            className={classNames(
              "hidden group-hover:block w-6 h-6 p-1 shrink-0 rounded-lg bg-gray-200 dark:bg-zinc-900 hover:opacity-80",
              alwaysShowLink && "!block",
            )}
            to={getShortcutUrl(shortcut.name)}
            target="_blank"
            onClick={(e) => e.stopPropagation()}
          >
            <Icon.ArrowUpRight className="w-3 h-auto text-gray-400 shrink-0" />
          </Link>
          {showActions && (
            <div className="ml-1" onClick={(e) => e.stopPropagation()}>
              <ShortcutActionsDropdown shortcut={shortcut} />
            </div>
          )}
        </div>
      </div>

      {/* Second row: shortcut link */}
      <div className="w-full mt-1 pl-7">
        <a
          className="text-xs text-blue-600 dark:text-blue-400 hover:underline truncate block"
          href={shortcutLink}
          target="_blank"
          onClick={(e) => e.stopPropagation()}
        >
          {shortcutLink}
        </a>
      </div>
    </div>
  );
};

export default ShortcutView;
