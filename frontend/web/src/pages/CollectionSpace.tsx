import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { useParams } from "react-router-dom";
import CollectionView from "@/components/CollectionView";
import Icon from "@/components/Icon";
import { useCollectionStore, useShortcutStore } from "@/stores";
import { Collection } from "@/types/proto/api/v1/collection_service";

const CollectionSpace = () => {
  const params = useParams();
  const collectionName = params.collectionName;
  const collectionStore = useCollectionStore();
  const shortcutStore = useShortcutStore();
  const [collection, setCollection] = useState<Collection>();
  const [loading, setLoading] = useState(true);

  if (!collectionName) {
    return null;
  }

  useEffect(() => {
    (async () => {
      try {
        setLoading(true);
        await shortcutStore.fetchShortcutList();
        const collection = await collectionStore.fetchCollectionByName(collectionName);
        setCollection(collection);
        document.title = `${collection.title} - Monotreme`;
      } catch (error: any) {
        console.error(error);
        toast.error(error.details);
      } finally {
        setLoading(false);
      }
    })();
  }, [collectionName]);

  if (loading) {
    return (
      <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 pt-4 pb-6 flex flex-col justify-center items-center py-16">
        <Icon.Loader className="w-8 h-auto animate-spin text-gray-400" />
        <p className="mt-4 text-gray-500">Loading collection...</p>
      </div>
    );
  }

  if (!collection) {
    return (
      <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 pt-4 pb-6 flex flex-col justify-center items-center py-16">
        <Icon.PackageOpen size={64} strokeWidth={1} className="text-gray-400" />
        <p className="mt-4 text-gray-500">Collection not found.</p>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 pt-4 pb-6 flex flex-col justify-start items-start">
      <CollectionView collection={collection} />
    </div>
  );
};

export default CollectionSpace;
