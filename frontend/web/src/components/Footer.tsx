import Icon from "./Icon";

const Footer: React.FC = () => {
  return (
    <footer className="w-full py-6 px-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-zinc-800">
      <div className="max-w-7xl mx-auto flex flex-col sm:flex-row justify-center items-center space-y-2 sm:space-y-0 sm:space-x-2 text-sm text-gray-600 dark:text-gray-400">
        <span>Built with love.</span>
        <span className="hidden sm:inline">•</span>
        <a
          href="https://github.com/bshort/monotreme"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center space-x-1 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
        >

          <span>Download the source code here</span>
        </a>
      </div>
      <div className="max-w-7xl mx-auto flex flex-col sm:flex-row justify-center items-center space-y-2 sm:space-y-0 sm:space-x-2 text-sm text-gray-600 dark:text-gray-400">
         <span>
           <Icon.Github className="size-6" strokeWidth={1.5}/> 
         </span>
         <span>
           https://github.com/bshort/monotreme
         </span>
      </div>
    </footer>
  );
};

export default Footer;