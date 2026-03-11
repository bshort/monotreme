// Utility functions for fetching URL metadata

const CORS_PROXY = "https://corsproxy.io/?";

/**
 * Fetches the title of a webpage from the given URL
 * @param url - The URL to fetch the title from
 * @returns Promise that resolves to the page title or null if failed
 */
export const fetchPageTitle = async (url: string): Promise<string | null> => {
  try {
    // Validate URL
    const urlObj = new URL(url);
    if (!['http:', 'https:'].includes(urlObj.protocol)) {
      return null;
    }

    // Create timeout controller
    const timeoutController = new AbortController();
    const timeoutId = setTimeout(() => timeoutController.abort(), 10000); // 10 seconds

    // Try to fetch with CORS proxy first
    let response: Response;
    try {
      response = await fetch(`${CORS_PROXY}${encodeURIComponent(url)}`, {
        method: 'GET',
        headers: {
          'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
        },
        signal: timeoutController.signal,
      });
    } catch (corsError) {
      // Fallback: try direct fetch (may fail due to CORS)
      try {
        const directController = new AbortController();
        const directTimeoutId = setTimeout(() => directController.abort(), 10000);

        response = await fetch(url, {
          method: 'GET',
          mode: 'cors',
          headers: {
            'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
          },
          signal: directController.signal,
        });

        clearTimeout(directTimeoutId);
      } catch (directError) {
        console.warn('Failed to fetch page title:', directError);
        clearTimeout(timeoutId);
        return null;
      }
    }

    clearTimeout(timeoutId);

    if (!response.ok) {
      return null;
    }

    const html = await response.text();

    // Extract title from HTML
    const titleMatch = html.match(/<title[^>]*>([^<]*)<\/title>/i);
    if (titleMatch && titleMatch[1]) {
      // Decode HTML entities and trim whitespace
      const title = titleMatch[1]
        .replace(/&amp;/g, '&')
        .replace(/&lt;/g, '<')
        .replace(/&gt;/g, '>')
        .replace(/&quot;/g, '"')
        .replace(/&#x27;/g, "'")
        .replace(/&#x2F;/g, '/')
        .replace(/&nbsp;/g, ' ')
        .trim();

      return title || null;
    }

    return null;
  } catch (error) {
    console.warn('Error fetching page title:', error);
    return null;
  }
};

/**
 * Generates a URL-friendly name from a title
 * @param title - The title to convert
 * @returns URL-friendly name (slug)
 */
export const generateUrlFriendlyName = (title: string): string => {
  if (!title || typeof title !== 'string') {
    return '';
  }

  return title
    .toLowerCase()
    // Remove special characters and replace with spaces
    .replace(/[^\w\s-]/g, ' ')
    // Replace multiple spaces/hyphens with single hyphen
    .replace(/[\s_-]+/g, '-')
    // Remove leading/trailing hyphens
    .replace(/^-+|-+$/g, '')
    // Limit length to reasonable size
    .substring(0, 50)
    // Remove trailing hyphen if substring created one
    .replace(/-+$/, '');
};

/**
 * Debounce function to limit API calls
 * @param func - Function to debounce
 * @param wait - Wait time in milliseconds
 * @returns Debounced function
 */
export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait: number
): ((...args: Parameters<T>) => void) => {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
};