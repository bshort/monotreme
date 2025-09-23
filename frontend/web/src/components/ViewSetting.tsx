import { Divider, Option, Select } from "@mui/joy";
import { useTranslation } from "react-i18next";
import { useViewStore } from "@/stores";
import Icon from "./Icon";
import Dropdown from "./common/Dropdown";

const ViewSetting = () => {
  const { t } = useTranslation();
  const viewStore = useViewStore();
  const order = viewStore.getOrder();
  const { field, direction } = order;
  const displayStyle = viewStore.displayStyle || "full";

  return (
    <Dropdown
      trigger={
        <button>
          <Icon.Settings2 className="w-4 h-auto text-gray-500" />
        </button>
      }
      actionsClassName="!mt-3 !right-[unset] -left-24 -ml-2"
      actions={
        <div className="w-52 p-2 gap-2 flex flex-col justify-start items-start" onClick={(e) => e.stopPropagation()}>
          <div className="w-full flex flex-col justify-start items-start gap-2">
            <span className="text-sm font-medium">{t("filter.display-mode")}</span>
            <div className="flex flex-col gap-1">
              <div className="flex flex-row justify-start items-center">
                <input
                  type="radio"
                  id="display-full"
                  name="display-mode"
                  className="mr-2"
                  checked={displayStyle === "full"}
                  onChange={() => viewStore.setDisplayStyle("full")}
                />
                <label htmlFor="display-full" className="text-sm cursor-pointer">
                  {t("filter.full-mode")}
                </label>
              </div>
              <div className="flex flex-row justify-start items-center">
                <input
                  type="radio"
                  id="display-compact"
                  name="display-mode"
                  className="mr-2"
                  checked={displayStyle === "compact"}
                  onChange={() => viewStore.setDisplayStyle("compact")}
                />
                <label htmlFor="display-compact" className="text-sm cursor-pointer">
                  {t("filter.compact-mode")}
                </label>
              </div>
              <div className="flex flex-row justify-start items-center">
                <input
                  type="radio"
                  id="display-list"
                  name="display-mode"
                  className="mr-2"
                  checked={displayStyle === "list"}
                  onChange={() => viewStore.setDisplayStyle("list")}
                />
                <label htmlFor="display-list" className="text-sm cursor-pointer">
                  {t("filter.list-mode")}
                </label>
              </div>
            </div>
          </div>
          <Divider className="!my-1" />
          <div className="w-full flex flex-row justify-between items-center">
            <span className="text-sm shrink-0 mr-2">{t("filter.order-by")}</span>
            <Select size="sm" value={field} onChange={(_, value) => viewStore.setOrder({ field: value as any })}>
              <Option value={"name"}>Name</Option>
              <Option value={"updatedTs"}>CreatedAt</Option>
              <Option value={"createdTs"}>UpdatedAt</Option>
              <Option value={"view"}>Visits</Option>
            </Select>
          </div>
          <div className="w-full flex flex-row justify-between items-center">
            <span className="text-sm shrink-0 mr-2">{t("filter.direction")}</span>
            <Select size="sm" value={direction} onChange={(_, value) => viewStore.setOrder({ direction: value as any })}>
              <Option value={"asc"}>ASC</Option>
              <Option value={"desc"}>DESC</Option>
            </Select>
          </div>
        </div>
      }
    ></Dropdown>
  );
};

export default ViewSetting;
