import { useTranslation } from "react-i18next";
import { useSearchParams, useNavigate } from "react-router-dom";
import { useViewStore } from "@/stores";
import Icon from "./Icon";
import VisibilityIcon from "./VisibilityIcon";

const FilterView = () => {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const viewStore = useViewStore();
  const filter = viewStore.filter;

  // Get current URL tags
  const urlTagsParam = searchParams.get("tags");
  const currentUrlTags = urlTagsParam ? urlTagsParam.split(",").map((tag) => tag.trim()) : [];

  const shouldShowFilters = filter.tag !== undefined || filter.visibility !== undefined || currentUrlTags.length > 0;

  const handleRemoveUrlTag = (tagToRemove: string) => {
    const newTags = currentUrlTags.filter((tag) => tag !== tagToRemove);
    if (newTags.length === 0) {
      navigate("/shortcuts");
    } else {
      navigate(`/shortcuts?tags=${newTags.map((t) => encodeURIComponent(t)).join(",")}`);
    }
  };

  if (!shouldShowFilters) {
    return <></>;
  }

  return (
    <div className="w-full flex flex-row justify-start items-center mb-4 pl-2">
      <span className="text-gray-400">Filters:</span>
      {currentUrlTags.map((tag) => (
        <button
          key={tag}
          className="ml-2 px-2 py-1 flex flex-row justify-center items-center bg-gray-100 rounded-full text-gray-500 text-sm hover:line-through"
          onClick={() => handleRemoveUrlTag(tag)}
        >
          <Icon.Tag className="w-4 h-auto mr-1" />
          <span className="max-w-[8rem] truncate">#{tag}</span>
          <Icon.X className="w-4 h-auto ml-1" />
        </button>
      ))}
      {filter.tag && (
        <button
          className="ml-2 px-2 py-1 flex flex-row justify-center items-center bg-gray-100 rounded-full text-gray-500 text-sm hover:line-through"
          onClick={() => viewStore.setFilter({ tag: undefined })}
        >
          <Icon.Tag className="w-4 h-auto mr-1" />
          <span className="max-w-[8rem] truncate">#{filter.tag}</span>
          <Icon.X className="w-4 h-auto ml-1" />
        </button>
      )}
      {filter.visibility && (
        <button
          className="ml-2 px-2 py-1 flex flex-row justify-center items-center bg-gray-100 rounded-full text-gray-500 text-sm hover:line-through"
          onClick={() => viewStore.setFilter({ visibility: undefined })}
        >
          <VisibilityIcon className="w-4 h-auto mr-1" visibility={filter.visibility} />
          {t(`shortcut.visibility.${filter.visibility.toLowerCase()}.self`)}
          <Icon.X className="w-4 h-auto ml-1" />
        </button>
      )}
    </div>
  );
};

export default FilterView;
