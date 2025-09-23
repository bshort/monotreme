import classNames from "classnames";
import { useTranslation } from "react-i18next";
import { Link, useSearchParams } from "react-router-dom";
import { useShortcutStore, useViewStore } from "@/stores";
import Icon from "./Icon";

const ShortcutsNavigator = () => {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const viewStore = useViewStore();
  const shortcutList = useShortcutStore().getShortcutList();
  const tags = shortcutList.map((shortcut) => shortcut.tags).flat();
  const currentTab = viewStore.filter.tab || `tab:all`;
  const sortedTagMap = sortTags(tags);

  // Get current URL tag filter
  const urlTagsParam = searchParams.get('tags');
  const currentUrlTag = urlTagsParam ? urlTagsParam : null;

  return (
    <div className="w-full flex flex-row justify-start items-center mb-4 gap-1 sm:flex-wrap overflow-x-auto no-scrollbar">
      <button
        className={classNames(
          "flex flex-row justify-center items-center px-2 leading-7 text-sm dark:text-gray-400 rounded-md",
          currentTab === "tab:all"
            ? "bg-blue-700 dark:bg-blue-800 text-white dark:text-gray-400 shadow"
            : "hover:bg-gray-200 dark:hover:bg-zinc-700",
        )}
        onClick={() => viewStore.setFilter({ tab: "tab:all" })}
      >
        <Icon.Earth className="w-4 h-auto mr-1" />
        <span className="font-normal">{t("filter.all")}</span>
      </button>
      <button
        className={classNames(
          "flex flex-row justify-center items-center px-2 leading-7 text-sm dark:text-gray-400 rounded-md",
          currentTab === "tab:mine"
            ? "bg-blue-700 dark:bg-blue-800 text-white dark:text-gray-400 shadow"
            : "hover:bg-gray-200 dark:hover:bg-zinc-700",
        )}
        onClick={() => viewStore.setFilter({ tab: "tab:mine" })}
      >
        <Icon.User className="w-4 h-auto mr-1" />
        <span className="font-normal">{t("filter.personal")}</span>
      </button>
      {Array.from(sortedTagMap.keys()).map((tag) => (
        <Link
          key={tag}
          to={`/shortcuts?tags=${encodeURIComponent(tag)}`}
          className={classNames(
            "flex flex-row justify-center items-center px-2 leading-7 text-sm dark:text-gray-400 rounded-md transition-colors",
            currentUrlTag === tag
              ? "bg-blue-700 dark:bg-blue-800 text-white dark:text-gray-400 shadow"
              : "hover:bg-gray-200 dark:hover:bg-zinc-700",
          )}
        >
          <Icon.Hash className="w-4 h-auto mr-0.5" />
          <span className="max-w-[8rem] truncate font-normal">{tag}</span>
        </Link>
      ))}
    </div>
  );
};

const sortTags = (tags: string[]): Map<string, number> => {
  const map = new Map<string, number>();
  for (const tag of tags) {
    const cleanTag = tag.trim().replace(/,$/, ''); // Remove trailing comma
    if (cleanTag) {
      const count = map.get(cleanTag) || 0;
      map.set(cleanTag, count + 1);
    }
  }
  const sortedMap = new Map([...map.entries()].sort((a, b) => b[1] - a[1]));
  return sortedMap;
};

export default ShortcutsNavigator;
