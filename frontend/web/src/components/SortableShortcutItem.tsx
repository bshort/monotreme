import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";

interface Props {
  shortcut: Shortcut;
  ShortcutItemView: React.ComponentType<any>;
  onClick: () => void;
}

const SortableShortcutItem: React.FC<Props> = ({ shortcut, ShortcutItemView, onClick }) => {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: shortcut.id,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
    cursor: isDragging ? "grabbing" : "grab",
  };

  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
      <ShortcutItemView shortcut={shortcut} showActions={true} onClick={onClick} />
    </div>
  );
};

export default SortableShortcutItem;
