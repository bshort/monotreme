const BrowserIntegrationSection: React.FC = () => {
  const bookmarkletCode = `javascript:(function(){try{const title=encodeURIComponent(document.title||'Untitled');const url=encodeURIComponent(window.location.href);window.open('${window.location.origin}/quick-save?url='+url+'&title='+title,'_blank','width=500,height=600,left='+(screen.width/2-250)+',top='+(screen.height/2-300)+',scrollbars=yes,resizable=yes');}catch(error){console.log('Mon.otre.me bookmarklet error:',error);}})();`;

  return (
    <div className="w-full flex flex-col sm:flex-row justify-start items-start gap-4 sm:gap-x-16">
      <p className="sm:w-1/4 text-2xl shrink-0 font-semibold text-gray-900 dark:text-gray-500">Browser Integration</p>
      <div className="w-full sm:w-auto grow flex flex-col justify-start items-start gap-4">
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
  );
};

export default BrowserIntegrationSection;
