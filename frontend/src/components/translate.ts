export type Language =
  | "zh"
  | "en"
  | "ja"
  | "ko"
  | "fr"
  | "de"
  | "es"
  | "it"
  | "ru"
  | "pt"
  | "vi"
  | "id"
  | "th"
  | "ms";

interface TranslateTemplate {
  header: {
    fn: string;
    client_key: string;
  };
  type: string;
  model_category: string;
  source: {
    lang: Language;
    text_list: string[];
  };
  target: {
    lang: Language;
  };
}

export interface TranslateOptions {
  from?: Language;
  to?: Language;
}

function getTranslateTemplate(
  textList: string[],
  from: Language = "zh",
  to: Language = "en"
): TranslateTemplate {
  return {
    header: {
      fn: "auto_translation",
      client_key:
        "browser-firefox-110.0.0-Windows 10-942844b6-8ddf-4e41-8851-b590f5022200-1685073533939",
    },
    type: "plain",
    model_category: "normal",
    source: {
      lang: from,
      text_list: textList,
    },
    target: {
      lang: to,
    },
  };
}

export async function translateAPI(
  text: string,
  options: TranslateOptions = {}
): Promise<string> {
  const { from = "zh", to = "en" } = options;
  let req = text;

  const data = JSON.stringify(getTranslateTemplate([req], from, to));
  const myHeaders = new Headers();
  myHeaders.append("Origin", "https://transmart.qq.com");
  myHeaders.append("Referer", "https://transmart.qq.com");
  myHeaders.append("X-Requested-With", "XMLHttpRequest");
  myHeaders.append(
    "User-Agent",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
  );
  myHeaders.append(
    "Cookie",
    "client_key=browser-firefox-110.0.0-Windows 10-942844b6-8ddf-4e41-8851-b590f5022201-" +
      parseInt(new Date().getTime().toString())
  );
  myHeaders.append("Content-Type", "application/json");

  const requestOptions: RequestInit = {
    method: "POST",
    headers: myHeaders,
    body: data,
    redirect: "follow" as RequestRedirect,
  };

  try {
    const resp = await (
      await fetch("https://transmart.qq.com/api/imt", requestOptions)
    ).json();
    return await resp.auto_translation[0];
  } catch (error) {
    return text;
  }
}
