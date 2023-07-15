import { useCallback, useEffect, useMemo, useState } from "react";
import { useLocation } from "react-router-dom";
import { Context } from "../app/context";
import ImagePreview from "./ImagePreview";

const DownloadPage: React.FC<{ context: Context }> = ({ context }) => {
  const location = useLocation();

  const [urls, setUrls] = useState<string[]>([]);
  const [asyncError, setAsyncError] = useState<unknown>();

  const key = useMemo(() => {
    const searchParams = new URLSearchParams(location.search);
    return searchParams.get('key');
  }, [location])

  const downloadAsync = useCallback(async (key: string) => {
    try {
      const info = await context.backend.download(key);
      const link = document.createElement('a');
      link.href = info.url;
      link.download = info.name;
      link.click();
    } catch (err) {
      setAsyncError(err)
    }
  }, []);

  useEffect(() => {
    if (!key) { return }

    const func = async () => {
      try {
        setUrls(await context.backend.getThumbnailUrl(key));
      } catch (err) {
        setAsyncError(err)
      }
    }
    func()
  }, [key])

  if (!key) {
    return <p>key is not specified</p>
  }

  return (
    <div>
      {asyncError ? <p>{`${asyncError}`}</p> : null}
      <button onClick={() => downloadAsync(key)}>Download</button>
      {
        urls.map((u, idx) => (<div key={idx}><ImagePreview src={u}/></div>))
      }
    </div>
  )
}

export default DownloadPage;
