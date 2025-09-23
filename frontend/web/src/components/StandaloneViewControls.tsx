import { Option, Select } from "@mui/joy";
import { useTranslation } from "react-i18next";
import { useShortcutStore, useViewStore } from "@/stores";
import { getAllUniqueTags } from "@/stores/view";

const StandaloneViewControls = () => {
  const { t } = useTranslation();
  const viewStore = useViewStore();
  const shortcutStore = useShortcutStore();
  const order = viewStore.getOrder();
  const { field, direction } = order;
  const displayStyle = viewStore.displayStyle || "full";
  const shortcutList = shortcutStore.getShortcutList();
  const allTags = getAllUniqueTags(shortcutList);
  const currentTagFilter = viewStore.filter.tag || "";

  return (
    <div className="flex flex-row justify-start items-center gap-4">
      {/* Display Mode Controls */}
      <div className="flex flex-row justify-start items-center gap-2">
        <span className="text-sm font-medium text-gray-600 dark:text-gray-400 whitespace-nowrap">
          {t("filter.display-mode")}:
        </span>
        <div className="flex flex-row justify-start items-center gap-3">
          <label className="flex flex-row justify-start items-center cursor-pointer">
            <input
              type="radio"
              name="display-mode"
              className="mr-1.5 accent-blue-600"
              checked={displayStyle === "full"}
              onChange={() => viewStore.setDisplayStyle("full")}
            />
            <span className="text-sm text-gray-700 dark:text-gray-300">
              {t("filter.full-mode")}
            </span>
          </label>
          <label className="flex flex-row justify-start items-center cursor-pointer">
            <input
              type="radio"
              name="display-mode"
              className="mr-1.5 accent-blue-600"
              checked={displayStyle === "compact"}
              onChange={() => viewStore.setDisplayStyle("compact")}
            />
            <span className="text-sm text-gray-700 dark:text-gray-300">
              {t("filter.compact-mode")}
            </span>
          </label>
          <label className="flex flex-row justify-start items-center cursor-pointer">
            <input
              type="radio"
              name="display-mode"
              className="mr-1.5 accent-blue-600"
              checked={displayStyle === "list"}
              onChange={() => viewStore.setDisplayStyle("list")}
            />
            <span className="text-sm text-gray-700 dark:text-gray-300">
              {t("filter.list-mode")}
            </span>
          </label>
        </div>
      </div>

      {/* Tags Filter Control */}
      <div className="flex flex-row justify-start items-center gap-2">
        <span className="text-sm font-medium text-gray-600 dark:text-gray-400 whitespace-nowrap">
          {t("filter.tag")}:
        </span>
        <Select
          size="sm"
          value={currentTagFilter}
          onChange={(_, value) => viewStore.setFilter({ tag: value || undefined })}
          className="min-w-[120px]"
          placeholder={t("filter.all-tags")}
        >
          <Option value="">
            {t("filter.all-tags")}
          </Option>
          {allTags.map((tag) => (
            <Option key={tag} value={tag}>
              #{tag}
            </Option>
          ))}
        </Select>
      </div>

      {/* Order By Control */}
      <div className="flex flex-row justify-start items-center gap-2">
        <span className="text-sm font-medium text-gray-600 dark:text-gray-400 whitespace-nowrap">
          {t("filter.order-by")}:
        </span>
        <Select
          size="sm"
          value={field}
          onChange={(_, value) => viewStore.setOrder({ field: value as any })}
          className="min-w-[100px]"
        >
          <Option value={"name"}>Name</Option>
          <Option value={"updatedTs"}>CreatedAt</Option>
          <Option value={"createdTs"}>UpdatedAt</Option>
          <Option value={"view"}>Visits</Option>
        </Select>
      </div>

      {/* Direction Control */}
      <div className="flex flex-row justify-start items-center gap-2">
        <span className="text-sm font-medium text-gray-600 dark:text-gray-400 whitespace-nowrap">
          {t("filter.direction")}:
        </span>
        <Select
          size="sm"
          value={direction}
          onChange={(_, value) => viewStore.setOrder({ direction: value as any })}
          className="min-w-[80px]"
        >
          <Option value={"asc"}>ASC</Option>
          <Option value={"desc"}>DESC</Option>
        </Select>
      </div>
    </div>
  );
};

export default StandaloneViewControls;