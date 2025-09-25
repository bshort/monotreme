import { Button, Option, Select, Switch } from "@mui/joy";
import { useTranslation } from "react-i18next";
import BetaBadge from "@/components/BetaBadge";
import { useUserStore } from "@/stores";
import { User } from "@/types/proto/api/v1/user_service";
import { Visibility } from "@/types/proto/api/v1/common";

const PreferenceSection: React.FC = () => {
  const { t } = useTranslation();
  const userStore = useUserStore();
  const currentUser = userStore.getCurrentUser();
  const language = currentUser.locale || "EN";
  const colorTheme = currentUser.colorTheme || "SYSTEM";
  const defaultVisibility = currentUser.defaultVisibility || "WORKSPACE";
  const autoGenerateTitle = currentUser.autoGenerateTitle ?? true;
  const autoGenerateIcon = currentUser.autoGenerateIcon ?? true;
  const autoGenerateName = currentUser.autoGenerateName ?? true;

  const bookmarkletCode = `javascript:(function(){try{const title=encodeURIComponent(document.title||'Untitled');const url=encodeURIComponent(window.location.href);window.open('${window.location.origin}/quick-save?url='+url+'&title='+title,'_blank','width=500,height=600,left='+(screen.width/2-250)+',top='+(screen.height/2-300)+',scrollbars=yes,resizable=yes');}catch(error){console.log('Mon.otre.me bookmarklet error:',error);}})();`;

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

  const handleSelectLanguage = async (locale: string) => {
    await userStore.patchUser(
      {
        ...currentUser,
        locale: locale,
      },
      ["locale"],
    );
  };

  const handleSelectColorTheme = async (colorTheme: string) => {
    await userStore.patchUser(
      {
        ...currentUser,
        colorTheme: colorTheme,
      },
      ["colorTheme"],
    );
  };

  const handleSelectDefaultVisibility = async (visibility: string) => {
    try {
      await userStore.patchUser(
        {
          ...currentUser,
          defaultVisibility: visibility,
        },
        ["defaultVisibility"],
      );
    } catch (error) {
      console.error('Failed to update default visibility setting:', error);
    }
  };

  const handleToggleAutoGenerateTitle = async (enabled: boolean) => {
    try {
      await userStore.patchUser(
        {
          ...currentUser,
          autoGenerateTitle: enabled,
        },
        ["autoGenerateTitle"],
      );
    } catch (error) {
      console.error('Failed to update auto-generate title setting:', error);
    }
  };

  const handleToggleAutoGenerateIcon = async (enabled: boolean) => {
    try {
      await userStore.patchUser(
        {
          ...currentUser,
          autoGenerateIcon: enabled,
        },
        ["autoGenerateIcon"],
      );
    } catch (error) {
      console.error('Failed to update auto-generate icon setting:', error);
    }
  };

  const handleToggleAutoGenerateName = async (enabled: boolean) => {
    try {
      await userStore.patchUser(
        {
          ...currentUser,
          autoGenerateName: enabled,
        },
        ["autoGenerateName"],
      );
    } catch (error) {
      console.error('Failed to update auto-generate name setting:', error);
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
          <Select defaultValue={colorTheme} onChange={(_, value) => handleSelectColorTheme(value as string)}>
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
          <Select defaultValue={language} onChange={(_, value) => handleSelectLanguage(value as string)}>
            {languageOptions.map((option) => {
              return (
                <Option key={option.value} value={option.value}>
                  {option.label}
                </Option>
              );
            })}
          </Select>
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

        {/* Browser Integration */}
        <div className="w-full border-t pt-4 mt-4 dark:border-zinc-700">
          <h4 className="text-lg font-semibold mb-3 dark:text-gray-300">Browser Integration</h4>

          <div className="w-full flex flex-col gap-3">
            <div className="flex flex-col">
              <span className="dark:text-gray-400 mb-1">Quick Save Bookmarklet</span>
              <span className="text-sm text-gray-500 dark:text-gray-600 mb-3">
                Drag this button to your bookmarks bar, then click it on any page to quickly save it to Mon.otre.me
              </span>
            </div>
            <div className="flex items-center gap-3">
              <a
                href={bookmarkletCode}
                className="inline-block px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg cursor-grab active:cursor-grabbing transition-colors duration-200 select-none no-underline font-medium"
                draggable={true}
                onClick={(e) => e.preventDefault()}
                title="Drag this to your bookmarks bar"
              >
                📚 Save to Mon.otre.me
              </a>
              <span className="text-xs text-gray-500 dark:text-gray-600">
                ← Drag this to your bookmarks bar
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PreferenceSection;
