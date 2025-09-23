import { Alert, Button, Card, Typography, Input, CircularProgress, Divider } from "@mui/joy";
import { useState } from "react";
import { useEffect } from "react";
import Icon from "@/components/Icon";
import useCollectionStore from "@/stores/collection";
import { useUserStore } from "@/stores";
import { Role } from "@/types/proto/api/v1/user_service";
import { ImportBookmarksResponse } from "@/types/proto/api/v1/collection_service";

const BookmarkImport = () => {
  const [file, setFile] = useState<File | null>(null);
  const [isImporting, setIsImporting] = useState(false);
  const [importResult, setImportResult] = useState<ImportBookmarksResponse | null>(null);
  const [error, setError] = useState<string>("");

  const collectionStore = useCollectionStore();
  const currentUser = useUserStore().getCurrentUser();
  const isAdmin = currentUser.role === Role.ADMIN;

  useEffect(() => {
    if (!isAdmin) {
      window.location.href = "/";
    }
  }, []);

  if (!isAdmin) {
    return null;
  }

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = event.target.files?.[0];
    if (selectedFile) {
      setFile(selectedFile);
      setError("");
      setImportResult(null);
    }
  };

  const handleImport = async () => {
    if (!file) {
      setError("Please select a bookmark file first");
      return;
    }

    setIsImporting(true);
    setError("");
    setImportResult(null);

    try {
      // Read the file as text
      const fileContent = await readFileAsText(file);

      // Call the import API
      const result = await collectionStore.importBookmarks(fileContent);

      setImportResult(result);
    } catch (err: any) {
      setError(err.message || "Failed to import bookmarks");
    } finally {
      setIsImporting(false);
    }
  };

  const readFileAsText = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => resolve(e.target?.result as string);
      reader.onerror = (e) => reject(new Error("Failed to read file"));
      reader.readAsText(file);
    });
  };

  const resetForm = () => {
    setFile(null);
    setError("");
    setImportResult(null);
    // Reset the file input
    const fileInput = document.getElementById("file-input") as HTMLInputElement;
    if (fileInput) {
      fileInput.value = "";
    }
  };

  return (
    <div className="mx-auto max-w-4xl w-full px-4 sm:px-6 md:px-8 py-6 flex flex-col justify-start items-start gap-y-8">
      <Alert variant="soft" color="warning" startDecorator={<Icon.Info />}>
        You can import bookmark files because you are an Admin.
      </Alert>

      <div className="w-full space-y-6">
        <div>
          <Typography level="h2" className="mb-4">
            Import Bookmarks
          </Typography>
          <Typography level="body-md" className="text-gray-600 dark:text-gray-400 mb-6">
            Upload an HTML bookmark file exported from your browser (Chrome, Firefox, Safari, Edge). The system will automatically detect the format and create collections for each bookmark folder and shortcuts for each bookmark.
          </Typography>
        </div>

        <Card className="w-full p-6">
          <div className="space-y-4">
            <div>
              <Typography level="title-md" className="mb-2">
                Select Bookmark File
              </Typography>
              <input
                id="file-input"
                type="file"
                accept=".html,.htm,text/html"
                onChange={handleFileChange}
                className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-medium file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
              />
              {file && (
                <Typography level="body-sm" className="mt-2 text-gray-600">
                  Selected: {file.name} ({Math.round(file.size / 1024)} KB)
                </Typography>
              )}
            </div>

            <div className="flex gap-3">
              <Button
                variant="solid"
                color="primary"
                onClick={handleImport}
                disabled={!file || isImporting}
                startDecorator={isImporting ? <CircularProgress size="sm" /> : <Icon.Upload />}
              >
                {isImporting ? "Importing..." : "Import Bookmarks"}
              </Button>

              {(file || importResult) && (
                <Button
                  variant="outlined"
                  color="neutral"
                  onClick={resetForm}
                  disabled={isImporting}
                >
                  Clear
                </Button>
              )}
            </div>
          </div>
        </Card>

        {error && (
          <Alert variant="solid" color="danger" startDecorator={<Icon.AlertTriangle />}>
            {error}
          </Alert>
        )}

        {importResult && (
          <Card className="w-full p-6">
            <Alert variant="solid" color="success" startDecorator={<Icon.CheckCircle />} className="mb-4">
              Bookmarks imported successfully!
            </Alert>

            <div className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="text-center">
                  <Typography level="h3" className="text-blue-600">
                    {importResult.totalCollections}
                  </Typography>
                  <Typography level="body-sm" className="text-gray-600">
                    Collections Created
                  </Typography>
                </div>
                <div className="text-center">
                  <Typography level="h3" className="text-green-600">
                    {importResult.totalShortcuts}
                  </Typography>
                  <Typography level="body-sm" className="text-gray-600">
                    Shortcuts Created
                  </Typography>
                </div>
                <div className="text-center">
                  <Typography level="h3" className="text-purple-600">
                    {importResult.collections.length}
                  </Typography>
                  <Typography level="body-sm" className="text-gray-600">
                    Collections Returned
                  </Typography>
                </div>
              </div>

              <Divider />

              <div>
                <Typography level="title-md" className="mb-3">
                  Created Collections:
                </Typography>
                <div className="space-y-2">
                  {importResult.collections.map((collection) => (
                    <div key={collection.id} className="flex justify-between items-center p-3 bg-gray-50 dark:bg-gray-800 rounded-md">
                      <div>
                        <Typography level="title-sm">{collection.title}</Typography>
                        <Typography level="body-xs" className="text-gray-600">
                          {collection.shortcutIds.length} shortcuts
                        </Typography>
                      </div>
                      <Typography level="body-xs" className="text-gray-500">
                        ID: {collection.id}
                      </Typography>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </Card>
        )}

        <Card className="w-full p-6 bg-blue-50 dark:bg-blue-900/20">
          <Typography level="title-md" className="mb-3 flex items-center gap-2">
            <Icon.Info /> Supported Formats & How it Works
          </Typography>
          <div className="space-y-3">
            <div>
              <Typography level="body-sm" className="font-medium mb-1">
                Supported Browser Formats:
              </Typography>
              <Typography level="body-sm" className="ml-4">
                • <strong>Chrome/Chromium:</strong> Standard HTML bookmarks export<br/>
                • <strong>Firefox:</strong> HTML bookmarks with Firefox-specific attributes<br/>
                • <strong>Safari/Edge:</strong> Standard HTML format (Chrome-compatible)
              </Typography>
            </div>

            <div>
              <Typography level="body-sm" className="font-medium mb-1">
                Import Process:
              </Typography>
              <Typography level="body-sm" className="ml-4">
                • System automatically detects your browser format<br/>
                • Bookmark folders (H3 headings) become collections<br/>
                • Individual bookmarks become shortcuts within collections<br/>
                • System folders and invalid links are automatically filtered out<br/>
                • All imported items use workspace visibility by default
              </Typography>
            </div>

            <div>
              <Typography level="body-sm" className="font-medium mb-1">
                How to Export Bookmarks:
              </Typography>
              <Typography level="body-sm" className="ml-4">
                • <strong>Chrome:</strong> Settings → Bookmarks → Bookmark Manager → Export bookmarks<br/>
                • <strong>Firefox:</strong> Library → Bookmarks → Show All Bookmarks → Import/Export → Export<br/>
                • <strong>Safari:</strong> File → Export Bookmarks<br/>
                • <strong>Edge:</strong> Settings → Profiles → Import/Export → Export to file
              </Typography>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default BookmarkImport;