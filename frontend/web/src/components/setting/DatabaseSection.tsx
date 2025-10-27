import { Button, Checkbox, Chip, Divider, Input, Radio, RadioGroup, Table, Typography } from "@mui/joy";
import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { workspaceServiceClient } from "@/grpcweb";
import Icon from "../Icon";

interface DatabaseStats {
  users: number;
  shortcuts: number;
  collections: number;
  tags: number;
  bookmarkTags: number;
  friendships: number;
  followings: number;
  activities: number;
  invitations: number;
}

interface EntityType {
  value: string;
  label: string;
  enabled: boolean;
}

type ImportMode = "overwrite" | "new-only" | "wipe-and-import";

const DatabaseSection = () => {
  const [importEntities, setImportEntities] = useState<EntityType[]>([
    { value: "users", label: "Users", enabled: false },
    { value: "shortcuts", label: "Shortcuts", enabled: false },
    { value: "collections", label: "Collections", enabled: false },
    { value: "tags", label: "Tags", enabled: false },
    { value: "friendships", label: "Friendships", enabled: false },
    { value: "followings", label: "Following", enabled: false },
  ]);

  const [exportEntities, setExportEntities] = useState<EntityType[]>([
    { value: "users", label: "Users", enabled: false },
    { value: "shortcuts", label: "Shortcuts", enabled: false },
    { value: "collections", label: "Collections", enabled: false },
    { value: "tags", label: "Tags", enabled: false },
    { value: "friendships", label: "Friendships", enabled: false },
    { value: "followings", label: "Following", enabled: false },
  ]);

  const [importFile, setImportFile] = useState<File | null>(null);
  const [importMode, setImportMode] = useState<ImportMode>("new-only");
  const [isImporting, setIsImporting] = useState(false);
  const [isExporting, setIsExporting] = useState(false);
  const [isLoadingStats, setIsLoadingStats] = useState(false);

  const [databaseStats, setDatabaseStats] = useState<DatabaseStats>({
    users: 0,
    shortcuts: 0,
    collections: 0,
    tags: 0,
    bookmarkTags: 0,
    friendships: 0,
    followings: 0,
    activities: 0,
    invitations: 0,
  });

  useEffect(() => {
    loadDatabaseStats();
  }, []);

  const loadDatabaseStats = async () => {
    setIsLoadingStats(true);
    try {
      const stats = await workspaceServiceClient.getDatabaseStats({});
      setDatabaseStats({
        users: stats.users,
        shortcuts: stats.shortcuts,
        collections: stats.collections,
        tags: stats.tags,
        bookmarkTags: stats.bookmarkTags,
        friendships: stats.friendships,
        followings: stats.followings,
        activities: stats.activities,
        invitations: stats.invitations,
      });
    } catch (error: any) {
      console.error("Failed to load database statistics:", error);
      toast.error("Failed to load database statistics");
    } finally {
      setIsLoadingStats(false);
    }
  };

  const toggleImportEntity = (value: string) => {
    setImportEntities(
      importEntities.map((entity) =>
        entity.value === value ? { ...entity, enabled: !entity.enabled } : entity
      )
    );
  };

  const toggleExportEntity = (value: string) => {
    setExportEntities(
      exportEntities.map((entity) =>
        entity.value === value ? { ...entity, enabled: !entity.enabled } : entity
      )
    );
  };

  const handleImportFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files && files.length > 0) {
      setImportFile(files[0]);
    }
  };

  const handleImport = async () => {
    const selectedEntities = importEntities.filter((e) => e.enabled);
    if (selectedEntities.length === 0) {
      toast.error("Please select at least one entity type to import");
      return;
    }

    if (!importFile) {
      toast.error("Please select a file to import");
      return;
    }

    // Show warning for destructive operations
    if (importMode === "wipe-and-import") {
      const confirmed = window.confirm(
        "⚠️ WARNING: This will DELETE ALL existing data in the selected tables before importing. This action cannot be undone. Are you sure you want to continue?"
      );
      if (!confirmed) {
        return;
      }
    }

    setIsImporting(true);
    try {
      // Read file content
      const fileContent = await importFile.text();

      // Call the backend API
      const response = await workspaceServiceClient.importDatabase({
        data: fileContent,
        entities: selectedEntities.map((e) => e.value),
        mode: importMode,
      });

      // Show success message with counts
      const counts = Object.entries(response.importedCounts)
        .map(([entity, count]) => `${entity}: ${count}`)
        .join(", ");

      toast.success(`Successfully imported: ${counts}`);

      // Show any warning messages
      if (response.messages && response.messages.length > 0) {
        response.messages.forEach((msg) => toast.error(msg));
      }

      setImportFile(null);
      setImportEntities(importEntities.map((e) => ({ ...e, enabled: false })));
      await loadDatabaseStats();
    } catch (error: any) {
      console.error("Import error:", error);
      toast.error(error.message || "Failed to import database");
    } finally {
      setIsImporting(false);
    }
  };

  const handleExport = async () => {
    const selectedEntities = exportEntities.filter((e) => e.enabled);
    if (selectedEntities.length === 0) {
      toast.error("Please select at least one entity type to export");
      return;
    }

    setIsExporting(true);
    try {
      // Call the backend API
      const response = await workspaceServiceClient.exportDatabase({
        entities: selectedEntities.map((e) => e.value),
      });

      // Create a blob from the JSON data
      const blob = new Blob([response.data], { type: "application/json" });
      const url = window.URL.createObjectURL(blob);

      // Create a temporary link and trigger download
      const a = document.createElement("a");
      a.href = url;
      a.download = response.filename;
      document.body.appendChild(a);
      a.click();

      // Cleanup
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);

      toast.success(`Exported ${selectedEntities.length} entity types successfully`);
      setExportEntities(exportEntities.map((e) => ({ ...e, enabled: false })));
    } catch (error: any) {
      console.error("Export error:", error);
      toast.error(error.message || "Failed to export database");
    } finally {
      setIsExporting(false);
    }
  };

  return (
    <div className="w-full flex flex-col sm:flex-row justify-start items-start gap-4 sm:gap-x-16">
      <p className="sm:w-1/4 text-2xl shrink-0 font-semibold text-gray-900 dark:text-gray-500">Database</p>
      <div className="w-full sm:w-auto grow flex flex-col justify-start items-start gap-8">

        {/* Database Import Section */}
        <div className="w-full flex flex-col gap-4">
          <div className="flex flex-row items-center gap-2">
            <Icon.Upload className="w-5 h-5 text-gray-600 dark:text-gray-400" />
            <Typography level="title-lg" className="font-semibold dark:text-gray-300">
              Import Database
            </Typography>
          </div>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Import data from a JSON file into the database. Select which entity types to import.
          </p>

          <div className="flex flex-col gap-2">
            <p className="text-sm font-medium dark:text-gray-400">Select entities to import:</p>
            <div className="flex flex-col gap-1 pl-2">
              {importEntities.map((entity) => (
                <Checkbox
                  key={entity.value}
                  label={entity.label}
                  checked={entity.enabled}
                  onChange={() => toggleImportEntity(entity.value)}
                />
              ))}
            </div>
          </div>

          <div className="flex flex-col gap-2">
            <p className="text-sm font-medium dark:text-gray-400">Import mode:</p>
            <RadioGroup
              value={importMode}
              onChange={(event) => setImportMode(event.target.value as ImportMode)}
              className="pl-2"
            >
              <Radio value="new-only" label="Import only new entries" />
              <Radio value="overwrite" label="Import and overwrite existing entries" />
              <Radio
                value="wipe-and-import"
                label={
                  <span className="flex items-center gap-1">
                    <Icon.AlertTriangle className="w-4 h-4 text-red-500" />
                    <span>Wipe database and import (DESTRUCTIVE)</span>
                  </span>
                }
              />
            </RadioGroup>
          </div>

          <div className="flex flex-col gap-2">
            <p className="text-sm font-medium dark:text-gray-400">Select file:</p>
            <div className="flex flex-row gap-2 items-center">
              <Input
                type="file"
                accept=".json,.sql"
                onChange={handleImportFileChange}
                className="flex-1"
              />
              {importFile && (
                <Chip color="success" size="sm">
                  {importFile.name}
                </Chip>
              )}
            </div>
          </div>

          <Button
            color="primary"
            loading={isImporting}
            disabled={isImporting || importEntities.filter((e) => e.enabled).length === 0 || !importFile}
            onClick={handleImport}
            startDecorator={<Icon.Upload className="w-4 h-4" />}
          >
            Import Data
          </Button>
        </div>

        <Divider />

        {/* Database Export Section */}
        <div className="w-full flex flex-col gap-4">
          <div className="flex flex-row items-center gap-2">
            <Icon.Download className="w-5 h-5 text-gray-600 dark:text-gray-400" />
            <Typography level="title-lg" className="font-semibold dark:text-gray-300">
              Export Database
            </Typography>
          </div>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Export database data to a JSON file. Select which entity types to export.
          </p>

          <div className="flex flex-col gap-2">
            <p className="text-sm font-medium dark:text-gray-400">Select entities to export:</p>
            <div className="flex flex-col gap-1 pl-2">
              {exportEntities.map((entity) => (
                <Checkbox
                  key={entity.value}
                  label={entity.label}
                  checked={entity.enabled}
                  onChange={() => toggleExportEntity(entity.value)}
                />
              ))}
            </div>
          </div>

          <Button
            color="primary"
            loading={isExporting}
            disabled={isExporting || exportEntities.filter((e) => e.enabled).length === 0}
            onClick={handleExport}
            startDecorator={<Icon.Download className="w-4 h-4" />}
          >
            Export Data
          </Button>
        </div>

        <Divider />

        {/* Database Report Section */}
        <div className="w-full flex flex-col gap-4">
          <div className="flex flex-row items-center gap-2">
            <Icon.BarChart className="w-5 h-5 text-gray-600 dark:text-gray-400" />
            <Typography level="title-lg" className="font-semibold dark:text-gray-300">
              Database Reports
            </Typography>
            <Button
              size="sm"
              variant="plain"
              onClick={loadDatabaseStats}
              loading={isLoadingStats}
              startDecorator={<Icon.RefreshCw className="w-4 h-4" />}
            >
              Refresh
            </Button>
          </div>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            View statistics and reports for each database table.
          </p>

          <Table
            borderAxis="both"
            stripe="odd"
            hoverRow
            className="mt-2"
          >
            <thead>
              <tr>
                <th style={{ width: "40%" }}>Table</th>
                <th style={{ width: "30%" }}>Records</th>
                <th style={{ width: "30%" }}>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>
                  <div className="flex items-center gap-2">
                    <Icon.Users className="w-4 h-4" />
                    <span className="font-medium">Users</span>
                  </div>
                </td>
                <td>{databaseStats.users.toLocaleString()}</td>
                <td>
                  <Chip size="sm" color={databaseStats.users > 0 ? "success" : "neutral"}>
                    {databaseStats.users > 0 ? "Active" : "Empty"}
                  </Chip>
                </td>
              </tr>
              <tr>
                <td>
                  <div className="flex items-center gap-2">
                    <Icon.Link className="w-4 h-4" />
                    <span className="font-medium">Shortcuts</span>
                  </div>
                </td>
                <td>{databaseStats.shortcuts.toLocaleString()}</td>
                <td>
                  <Chip size="sm" color={databaseStats.shortcuts > 0 ? "success" : "neutral"}>
                    {databaseStats.shortcuts > 0 ? "Active" : "Empty"}
                  </Chip>
                </td>
              </tr>
              <tr>
                <td>
                  <div className="flex items-center gap-2">
                    <Icon.Folder className="w-4 h-4" />
                    <span className="font-medium">Collections</span>
                  </div>
                </td>
                <td>{databaseStats.collections.toLocaleString()}</td>
                <td>
                  <Chip size="sm" color={databaseStats.collections > 0 ? "success" : "neutral"}>
                    {databaseStats.collections > 0 ? "Active" : "Empty"}
                  </Chip>
                </td>
              </tr>
              <tr>
                <td>
                  <div className="flex items-center gap-2">
                    <Icon.Hash className="w-4 h-4" />
                    <span className="font-medium">Tags</span>
                  </div>
                </td>
                <td>{databaseStats.tags.toLocaleString()}</td>
                <td>
                  <Chip size="sm" color={databaseStats.tags > 0 ? "success" : "neutral"}>
                    {databaseStats.tags > 0 ? "Active" : "Empty"}
                  </Chip>
                </td>
              </tr>
              <tr>
                <td>
                  <div className="flex items-center gap-2">
                    <Icon.Tag className="w-4 h-4" />
                    <span className="font-medium">Bookmark Tags</span>
                  </div>
                </td>
                <td>{databaseStats.bookmarkTags.toLocaleString()}</td>
                <td>
                  <Chip size="sm" color={databaseStats.bookmarkTags > 0 ? "success" : "neutral"}>
                    {databaseStats.bookmarkTags > 0 ? "Active" : "Empty"}
                  </Chip>
                </td>
              </tr>
              <tr>
                <td>
                  <div className="flex items-center gap-2">
                    <Icon.UserCheck className="w-4 h-4" />
                    <span className="font-medium">Friendships</span>
                  </div>
                </td>
                <td>{databaseStats.friendships.toLocaleString()}</td>
                <td>
                  <Chip size="sm" color={databaseStats.friendships > 0 ? "success" : "neutral"}>
                    {databaseStats.friendships > 0 ? "Active" : "Empty"}
                  </Chip>
                </td>
              </tr>
              <tr>
                <td>
                  <div className="flex items-center gap-2">
                    <Icon.UserPlus className="w-4 h-4" />
                    <span className="font-medium">Following</span>
                  </div>
                </td>
                <td>{databaseStats.followings.toLocaleString()}</td>
                <td>
                  <Chip size="sm" color={databaseStats.followings > 0 ? "success" : "neutral"}>
                    {databaseStats.followings > 0 ? "Active" : "Empty"}
                  </Chip>
                </td>
              </tr>
              <tr>
                <td>
                  <div className="flex items-center gap-2">
                    <Icon.Activity className="w-4 h-4" />
                    <span className="font-medium">Activities</span>
                  </div>
                </td>
                <td>{databaseStats.activities.toLocaleString()}</td>
                <td>
                  <Chip size="sm" color={databaseStats.activities > 0 ? "success" : "neutral"}>
                    {databaseStats.activities > 0 ? "Active" : "Empty"}
                  </Chip>
                </td>
              </tr>
              <tr>
                <td>
                  <div className="flex items-center gap-2">
                    <Icon.Mail className="w-4 h-4" />
                    <span className="font-medium">Invitations</span>
                  </div>
                </td>
                <td>{databaseStats.invitations.toLocaleString()}</td>
                <td>
                  <Chip size="sm" color={databaseStats.invitations > 0 ? "success" : "neutral"}>
                    {databaseStats.invitations > 0 ? "Active" : "Empty"}
                  </Chip>
                </td>
              </tr>
            </tbody>
          </Table>
        </div>
      </div>
    </div>
  );
};

export default DatabaseSection;
