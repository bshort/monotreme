import { Link } from "@mui/joy";
import { useTranslation } from "react-i18next";
import aboutContent from "@/content/about.json";

interface ContentSection {
  type: string;
  content?: string;
  label?: string;
  url?: string;
  text?: string;
}

interface AboutContent {
  title: string;
  description: string;
  sections: ContentSection[];
}

const About: React.FC = () => {
  const { t } = useTranslation();
  const content: AboutContent = aboutContent;

  const renderSection = (section: ContentSection, index: number) => {
    switch (section.type) {
      case "paragraph":
        return (
          <p key={index} className="text-lg">
            {section.content}
          </p>
        );
      case "link":
        return (
          <div key={index} className="mt-4">
            <span className="mr-2">{section.label}</span>
            <Link variant="plain" href={section.url} target="_blank">
              {section.text}
            </Link>
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <div className="w-full max-w-2xl mx-auto">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">{t("common.about")}</h1>
      </div>

      <div className="space-y-4 text-gray-700 dark:text-gray-300">
        {content.sections.map((section, index) => renderSection(section, index))}
      </div>
    </div>
  );
};

export default About;