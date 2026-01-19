import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { Collection } from "@/types/proto/api/v1/collection_service";
import CollectionView from "./CollectionView";

interface Props {
  collection: Collection;
  searchQuery?: string;
}

const SortableCollectionItem: React.FC<Props> = ({ collection, searchQuery }) => {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: collection.id,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const dragHandleProps = {
    ...attributes,
    ...listeners,
  };

  return (
    <div ref={setNodeRef} style={style} className="w-full">
      <CollectionView collection={collection} searchQuery={searchQuery} dragHandleProps={dragHandleProps} />
    </div>
  );
};

export default SortableCollectionItem;
