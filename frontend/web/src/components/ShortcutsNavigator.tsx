import classNames from "classnames";
import { useTranslation } from "react-i18next";
import { useSearchParams, useNavigate } from "react-router-dom";
import { useShortcutStore, useViewStore } from "@/stores";
import Icon from "./Icon";

const ShortcutsNavigator = () => {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const viewStore = useViewStore();
  const shortcutList = useShortcutStore().getShortcutList();
  const tags = shortcutList.map((shortcut) => shortcut.tags).flat();
  const currentTab = viewStore.filter.tab || `tab:all`;
  const sortedTagMap = sortTags(tags);

  // Get current URL tag filter
  const urlTagsParam = searchParams.get("tags");
  const currentUrlTags = urlTagsParam ? urlTagsParam.split(",").map((tag) => tag.trim()) : [];

  const handleTagClick = (tag: string, event: React.MouseEvent) => {
    event.preventDefault();

    if (currentUrlTags.includes(tag)) {
      // Tag is selected, unselect it
      const newTags = currentUrlTags.filter((t) => t !== tag);
      if (newTags.length === 0) {
        // No tags left, go to /shortcuts
        navigate("/shortcuts");
      } else {
        // Still have tags, update URL
        navigate(`/shortcuts?tags=${newTags.map((t) => encodeURIComponent(t)).join(",")}`);
      }
    } else {
      // Tag is not selected, add it to existing tags
      const newTags = [...currentUrlTags, tag];
      navigate(`/shortcuts?tags=${newTags.map((t) => encodeURIComponent(t)).join(",")}`);
    }
  };

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
        <button
          key={tag}
          onClick={(event) => handleTagClick(tag, event)}
          className={classNames(
            "flex flex-row justify-center items-center px-2 leading-7 text-sm dark:text-gray-400 rounded-md transition-colors",
            currentUrlTags.includes(tag)
              ? "bg-blue-700 dark:bg-blue-800 text-white dark:text-gray-400 shadow"
              : "hover:bg-gray-200 dark:hover:bg-zinc-700",
          )}
        >
          <Icon.Hash className="w-4 h-auto mr-0.5" />
          <span className="max-w-[8rem] truncate font-normal">{tag}</span>
        </button>
      ))}
    </div>
  );
};

const sortTags = (tags: string[]): Map<string, number> => {
  const map = new Map<string, number>();
  for (const tag of tags) {
    const cleanTag = tag.trim().replace(/,$/, ""); // Remove trailing comma
    if (cleanTag) {
      const count = map.get(cleanTag) || 0;
      map.set(cleanTag, count + 1);
    }
  }
  const sortedMap = new Map([...map.entries()].sort((a, b) => b[1] - a[1]));
  return sortedMap;
};

export default ShortcutsNavigator;
