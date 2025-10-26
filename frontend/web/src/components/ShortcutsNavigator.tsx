import classNames from "classnames";
import { useTranslation } from "react-i18next";
import { useViewStore } from "@/stores";
import Icon from "./Icon";

const ShortcutsNavigator = () => {
  const { t } = useTranslation();
  const viewStore = useViewStore();
  const currentTab = viewStore.filter.tab || `tab:all`;

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
        <span className="font-normal">My</span>
      </button>
    </div>
  );
};

export default ShortcutsNavigator;
