import classNames from "classnames";
import useNavigateTo from "@/hooks/useNavigateTo";
import { useViewStore, useWorkspaceStore } from "@/stores";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";
import { getShortcutUrl } from "@/utils/shortcut";
import ShortcutCard from "./ShortcutCard";
import ShortcutView from "./ShortcutView";
import ShortcutListView from "./ShortcutListView";

interface Props {
  shortcutList: Shortcut[];
}

const ShortcutsContainer: React.FC<Props> = (props: Props) => {
  const { shortcutList } = props;
  const navigateTo = useNavigateTo();
  const viewStore = useViewStore();
  const workspaceStore = useWorkspaceStore();
  const displayStyle = viewStore.displayStyle || "full";

  let ShortcutItemView = ShortcutCard;
  if (displayStyle === "compact") {
    ShortcutItemView = ShortcutView;
  } else if (displayStyle === "list") {
    ShortcutItemView = ShortcutListView;
  }

  const handleShortcutClick = (shortcut: Shortcut) => {
    // Use the server route to ensure visit tracking
    const shortcutUrl = getShortcutUrl(shortcut.name);
    window.open(shortcutUrl, '_blank');
  };

  let gridClasses = "w-full grid grid-cols-1 gap-3 sm:gap-4";
  if (displayStyle === "full") {
    gridClasses += " sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4";
  } else if (displayStyle === "compact") {
    gridClasses += " grid-cols-2 sm:grid-cols-4";
  } else if (displayStyle === "list") {
    gridClasses += " grid-cols-1";
  }

  return (
    <div className={gridClasses}>
      {shortcutList.map((shortcut) => {
        return <ShortcutItemView key={shortcut.id} shortcut={shortcut} showActions={true} onClick={() => handleShortcutClick(shortcut)} />;
      })}
    </div>
  );
};

export default ShortcutsContainer;
