import { DndContext, closestCenter, KeyboardSensor, PointerSensor, useSensor, useSensors, DragEndEvent } from "@dnd-kit/core";
import { arrayMove, SortableContext, sortableKeyboardCoordinates, verticalListSortingStrategy, rectSortingStrategy } from "@dnd-kit/sortable";
import classNames from "classnames";
import { useState, useEffect } from "react";
import useNavigateTo from "@/hooks/useNavigateTo";
import { useShortcutStore, useViewStore, useWorkspaceStore } from "@/stores";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";
import { getShortcutUrl } from "@/utils/shortcut";
import ShortcutCard from "./ShortcutCard";
import ShortcutView from "./ShortcutView";
import ShortcutListView from "./ShortcutListView";
import SortableShortcutItem from "./SortableShortcutItem";

interface Props {
  shortcutList: Shortcut[];
}

const ShortcutsContainer: React.FC<Props> = (props: Props) => {
  const { shortcutList: initialShortcutList } = props;
  const [shortcutList, setShortcutList] = useState(initialShortcutList);
  const navigateTo = useNavigateTo();
  const viewStore = useViewStore();
  const shortcutStore = useShortcutStore();
  const workspaceStore = useWorkspaceStore();
  const displayStyle = viewStore.displayStyle || "full";
  const order = viewStore.getOrder();
  const isCustomOrder = order.field === "userOrder";

  // Update local state when prop changes
  useEffect(() => {
    setShortcutList(initialShortcutList);
  }, [initialShortcutList]);

  let ShortcutItemView = ShortcutCard;
  if (displayStyle === "compact") {
    ShortcutItemView = ShortcutView;
  } else if (displayStyle === "list") {
    ShortcutItemView = ShortcutListView;
  }

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const handleShortcutClick = (shortcut: Shortcut) => {
    // Use the server route to ensure visit tracking
    const shortcutUrl = getShortcutUrl(shortcut.name);
    window.open(shortcutUrl, '_blank');
  };

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;

    if (!over || active.id === over.id) {
      return;
    }

    const oldIndex = shortcutList.findIndex((s) => s.id === active.id);
    const newIndex = shortcutList.findIndex((s) => s.id === over.id);

    const newShortcutList = arrayMove(shortcutList, oldIndex, newIndex);
    setShortcutList(newShortcutList);

    // Update userOrder for all shortcuts and save to backend
    const updatedShortcuts = newShortcutList.map((shortcut, index) => ({
      ...shortcut,
      userOrder: index,
    }));

    // Save all shortcuts with updated userOrder
    try {
      await Promise.all(
        updatedShortcuts.map((shortcut) =>
          shortcutStore.updateShortcut({
            ...shortcut,
          }, ["user_order"])
        )
      );
      // Refresh the shortcut list to get updated data
      await shortcutStore.fetchShortcutList();
    } catch (error) {
      console.error("Failed to update shortcut order:", error);
      // Revert to original order on error
      setShortcutList(initialShortcutList);
    }
  };

  let gridClasses = "w-full grid grid-cols-1 gap-3 sm:gap-4";
  if (displayStyle === "full") {
    gridClasses += " sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4";
  } else if (displayStyle === "compact") {
    gridClasses += " grid-cols-2 sm:grid-cols-4";
  } else if (displayStyle === "list") {
    gridClasses += " grid-cols-1";
  }

  // Only enable drag-and-drop when in custom order mode
  if (!isCustomOrder) {
    return (
      <div className={gridClasses}>
        {shortcutList.map((shortcut) => {
          return <ShortcutItemView key={shortcut.id} shortcut={shortcut} showActions={true} onClick={() => handleShortcutClick(shortcut)} />;
        })}
      </div>
    );
  }

  const sortingStrategy = displayStyle === "list" ? verticalListSortingStrategy : rectSortingStrategy;

  return (
    <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
      <SortableContext items={shortcutList.map((s) => s.id)} strategy={sortingStrategy}>
        <div className={gridClasses}>
          {shortcutList.map((shortcut) => {
            return (
              <SortableShortcutItem
                key={shortcut.id}
                shortcut={shortcut}
                ShortcutItemView={ShortcutItemView}
                onClick={() => handleShortcutClick(shortcut)}
              />
            );
          })}
        </div>
      </SortableContext>
    </DndContext>
  );
};

export default ShortcutsContainer;
