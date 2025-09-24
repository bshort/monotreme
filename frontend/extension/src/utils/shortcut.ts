import { Storage } from "@plasmohq/storage";

const storage = new Storage();

export const getShortcutUrl = async (shortcutName: string): Promise<string> => {
  // For the extension, we'll use a default prefix of 's' unless configured otherwise
  // This could be enhanced to read from storage or configuration
  const prefix = await getShortcutPrefix();
  return `/${prefix}/${shortcutName}`;
};

export const getShortcutPrefix = async (): Promise<string> => {
  // Try to get the prefix from storage, fallback to 's'
  const prefix = await storage.getItem<string>("shortcut_prefix");
  return prefix || "s";
};