import { create } from "zustand";
import { collectionServiceClient } from "@/grpcweb";
import { Collection, ImportBookmarksResponse } from "@/types/proto/api/v1/collection_service";

interface CollectionState {
  collectionMapById: Record<number, Collection>;
  fetchCollectionList: () => Promise<Collection[]>;
  getOrFetchCollectionById: (id: number) => Promise<Collection>;
  getCollectionById: (id: number) => Collection;
  getCollectionList: () => Collection[];
  fetchCollectionByName: (collectionName: string) => Promise<Collection>;
  createCollection: (collection: Collection) => Promise<Collection>;
  updateCollection: (collection: Partial<Collection>, updateMask: string[]) => Promise<Collection>;
  deleteCollection: (id: number) => Promise<void>;
  importBookmarks: (htmlContent: string) => Promise<ImportBookmarksResponse>;
}

const useCollectionStore = create<CollectionState>()((set, get) => ({
  collectionMapById: {},
  fetchCollectionList: async () => {
    const { collections } = await collectionServiceClient.listCollections({});
    const collectionMap = get().collectionMapById;
    collections.forEach((collection) => {
      collectionMap[collection.id] = collection;
    });
    set(collectionMap);
    return collections;
  },
  getOrFetchCollectionById: async (id: number) => {
    const collectionMap = get().collectionMapById;
    if (collectionMap[id]) {
      return collectionMap[id] as Collection;
    }

    const collection = await collectionServiceClient.getCollection({
      id: id,
    });
    collectionMap[id] = collection;
    set(collectionMap);
    return collection;
  },
  getCollectionById: (id: number) => {
    const collectionMap = get().collectionMapById;
    return collectionMap[id] as Collection;
  },
  getCollectionList: () => {
    return Object.values(get().collectionMapById);
  },
  fetchCollectionByName: async (collectionName: string) => {
    const collection = await collectionServiceClient.getCollectionByName({
      name: collectionName,
    });
    const collectionMap = get().collectionMapById;
    collectionMap[collection.id] = collection;
    set(collectionMap);
    return collection;
  },
  createCollection: async (collection: Collection) => {
    const createdCollection = await collectionServiceClient.createCollection({
      collection: collection,
    });
    const collectionMap = get().collectionMapById;
    collectionMap[createdCollection.id] = createdCollection;
    set(collectionMap);
    return createdCollection;
  },
  updateCollection: async (collection: Partial<Collection>, updateMask: string[]) => {
    const updatedCollection = await collectionServiceClient.updateCollection({
      collection: collection,
      updateMask: updateMask,
    });
    const collectionMap = get().collectionMapById;
    collectionMap[updatedCollection.id] = updatedCollection;
    set(collectionMap);
    return updatedCollection;
  },
  deleteCollection: async (id: number) => {
    await collectionServiceClient.deleteCollection({
      id,
    });
    const collectionMap = get().collectionMapById;
    delete collectionMap[id];
    set(collectionMap);
  },
  importBookmarks: async (htmlContent: string) => {
    const response = await collectionServiceClient.importBookmarks({
      htmlContent,
    });

    // Update the collection map with newly created collections
    const collectionMap = get().collectionMapById;
    response.collections.forEach((collection) => {
      collectionMap[collection.id] = collection;
    });
    set(collectionMap);

    return response;
  },
}));

export default useCollectionStore;
