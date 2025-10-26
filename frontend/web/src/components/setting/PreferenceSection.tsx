import { Button, Option, Select, Switch } from "@mui/joy";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import toast from "react-hot-toast";
import BetaBadge from "@/components/BetaBadge";
import { useUserStore } from "@/stores";
import { User } from "@/types/proto/api/v1/user_service";
import { Visibility } from "@/types/proto/api/v1/common";

const PreferenceSection: React.FC = () => {
  const { t } = useTranslation();
  const userStore = useUserStore();
  const currentUser = userStore.getCurrentUser();

  // Local state for pending changes
  const [language, setLanguage] = useState(currentUser.locale || "EN");
  const [colorTheme, setColorTheme] = useState(currentUser.colorTheme || "SYSTEM");
  const [defaultVisibility, setDefaultVisibility] = useState(currentUser.defaultVisibility || "WORKSPACE");
  const [autoGenerateTitle, setAutoGenerateTitle] = useState(currentUser.autoGenerateTitle ?? true);
  const [autoGenerateIcon, setAutoGenerateIcon] = useState(currentUser.autoGenerateIcon ?? true);
  const [autoGenerateName, setAutoGenerateName] = useState(currentUser.autoGenerateName ?? true);
  const [editModePreference, setEditModePreference] = useState(currentUser.editModePreference || "FLYOUT");
  const [hasUnsavedChanges, setHasUnsavedChanges] = useState(false);

  const languageOptions = [
    {
      value: "EN",
      label: "English",
    },
    {
      value: "ZH",
      label: "中文",
    },
    {
      value: "FR",
      label: "Français",
    },
    {
      value: "JA",
      label: "日本語",
    },
    {
      value: "TR",
      label: "Türkçe",
    },
    {
      value: "RU",
      label: "русский",
    },
    {
      value: "HU",
      label: "Magyar",
    },
  ];

  const colorThemeOptions = [
    {
      value: "SYSTEM",
      label: "System",
    },
    {
      value: "LIGHT",
      label: "Light",
    },
    {
      value: "DARK",
      label: "Dark",
    },
  ];

  const visibilityOptions = [
    {
      value: "PRIVATE",
      label: "Private",
    },
    {
      value: "WORKSPACE",
      label: "Workspace",
    },
    {
      value: "PUBLIC",
      label: "Public",
    },
  ];

  const editModeOptions = [
    {
      value: "FLYOUT",
      label: "Quick Edit Flyout",
    },
    {
      value: "FULL_PAGE",
      label: "Full Page Edit",
    },
  ];

  const handleSelectLanguage = (locale: string) => {
    setLanguage(locale);
    setHasUnsavedChanges(true);
  };

  const handleSelectColorTheme = (theme: string) => {
    setColorTheme(theme);
    setHasUnsavedChanges(true);
  };

  const handleSelectDefaultVisibility = (visibility: string) => {
    setDefaultVisibility(visibility);
    setHasUnsavedChanges(true);
  };

  const handleToggleAutoGenerateTitle = (enabled: boolean) => {
    setAutoGenerateTitle(enabled);
    setHasUnsavedChanges(true);
  };

  const handleToggleAutoGenerateIcon = (enabled: boolean) => {
    setAutoGenerateIcon(enabled);
    setHasUnsavedChanges(true);
  };

  const handleToggleAutoGenerateName = (enabled: boolean) => {
    setAutoGenerateName(enabled);
    setHasUnsavedChanges(true);
  };

  const handleSelectEditModePreference = (mode: string) => {
    setEditModePreference(mode);
    setHasUnsavedChanges(true);
  };

  const handleSave = async () => {
    try {
      await userStore.patchUser(
        {
          ...currentUser,
          locale: language,
          colorTheme: colorTheme,
          defaultVisibility: defaultVisibility,
          autoGenerateTitle: autoGenerateTitle,
          autoGenerateIcon: autoGenerateIcon,
          autoGenerateName: autoGenerateName,
          editModePreference: editModePreference,
        },
        [
          "locale",
          "colorTheme",
          "defaultVisibility",
          "autoGenerateTitle",
          "autoGenerateIcon",
          "autoGenerateName",
          "editModePreference",
        ],
      );
      setHasUnsavedChanges(false);
      toast.success("Preferences saved successfully");
    } catch (error) {
      console.error('Failed to save preferences:', error);
      toast.error("Failed to save preferences. Please try again.");
    }
  };

  return (
    <div className="w-full flex flex-col sm:flex-row justify-start items-start gap-4 sm:gap-x-16">
      <p className="sm:w-1/4 text-2xl shrink-0 font-semibold text-gray-900 dark:text-gray-500">{t("settings.preference.self")}</p>
      <div className="w-full sm:w-auto grow flex flex-col justify-start items-start gap-4">
        <div className="w-full flex flex-row justify-between items-center">
          <div className="flex flex-row justify-start items-center gap-x-1">
            <span className="dark:text-gray-400">{t("settings.preference.color-theme")}</span>
          </div>
          <Select value={colorTheme} onChange={(_, value) => handleSelectColorTheme(value as string)}>
            {colorThemeOptions.map((option) => {
              return (
                <Option key={option.value} value={option.value}>
                  {option.label}
                </Option>
              );
            })}
          </Select>
        </div>
        <div className="w-full flex flex-row justify-between items-center">
          <div className="flex flex-row justify-start items-center gap-x-1">
            <span className="dark:text-gray-400">{t("common.language")}</span>
            <BetaBadge />
          </div>
          <Select value={language} onChange={(_, value) => handleSelectLanguage(value as string)}>
            {languageOptions.map((option) => {
              return (
                <Option key={option.value} value={option.value}>
                  {option.label}
                </Option>
              );
            })}
          </Select>
        </div>

        {/* Shortcut Edit Preferences */}
        <div className="w-full border-t pt-4 mt-2 dark:border-zinc-700">
          <h4 className="text-lg font-semibold mb-3 dark:text-gray-300">Shortcut Edit Preferences</h4>

          <div className="w-full flex flex-row justify-between items-center">
            <div className="flex flex-col">
              <span className="dark:text-gray-400">Default Edit Mode</span>
              <span className="text-sm text-gray-500 dark:text-gray-600">Choose how to edit shortcuts by default</span>
            </div>
            <Select value={editModePreference} onChange={(_, value) => handleSelectEditModePreference(value as string)}>
              {editModeOptions.map((option) => {
                return (
                  <Option key={option.value} value={option.value}>
                    {option.label}
                  </Option>
                );
              })}
            </Select>
          </div>
        </div>

        {/* Shortcut Creation Preferences */}
        <div className="w-full border-t pt-4 mt-2 dark:border-zinc-700">
          <h4 className="text-lg font-semibold mb-3 dark:text-gray-300">Shortcut Creation Preferences</h4>

          <div className="w-full flex flex-row justify-between items-center mb-3">
            <span className="dark:text-gray-400">Default Visibility</span>
            <Select value={defaultVisibility} onChange={(_, value) => handleSelectDefaultVisibility(value as string)}>
              {visibilityOptions.map((option) => {
                return (
                  <Option key={option.value} value={option.value}>
                    {option.label}
                  </Option>
                );
              })}
            </Select>
          </div>

          <div className="w-full flex flex-row justify-between items-center mb-3">
            <div className="flex flex-col">
              <span className="dark:text-gray-400">Auto-generate Title</span>
              <span className="text-sm text-gray-500 dark:text-gray-600">Automatically fetch page title from URL</span>
            </div>
            <Switch
              checked={autoGenerateTitle}
              onChange={(event) => handleToggleAutoGenerateTitle(event.target.checked)}
            />
          </div>

          <div className="w-full flex flex-row justify-between items-center mb-3">
            <div className="flex flex-col">
              <span className="dark:text-gray-400">Auto-generate Icon</span>
              <span className="text-sm text-gray-500 dark:text-gray-600">Automatically fetch favicon from website</span>
            </div>
            <Switch
              checked={autoGenerateIcon}
              onChange={(event) => handleToggleAutoGenerateIcon(event.target.checked)}
            />
          </div>

          <div className="w-full flex flex-row justify-between items-center">
            <div className="flex flex-col">
              <span className="dark:text-gray-400">Auto-generate Shortcut</span>
              <span className="text-sm text-gray-500 dark:text-gray-600">Automatically create URL-friendly shortcut name from title</span>
            </div>
            <Switch
              checked={autoGenerateName}
              onChange={(event) => handleToggleAutoGenerateName(event.target.checked)}
            />
          </div>
        </div>

        {/* Save Button */}
        <div className="w-full border-t pt-4 mt-4 dark:border-zinc-700">
          <Button
            color="primary"
            size="lg"
            onClick={handleSave}
            disabled={!hasUnsavedChanges}
            sx={{
              fontWeight: 'bold',
              fontSize: '1.125rem',
              paddingX: '2rem',
              paddingY: '0.875rem',
              minWidth: '150px',
            }}
          >
            {t("common.save")}
          </Button>
          {hasUnsavedChanges && (
            <span className="ml-3 text-sm text-orange-600 dark:text-orange-400">
              You have unsaved changes
            </span>
          )}
        </div>
      </div>
    </div>
  );
};

export default PreferenceSection;
