import { useMemo } from "react";
import { useLocation } from "react-router-dom";
import { Context } from "../app/context";

const DownloadPage: React.FC<{ context: Context }> = ({ context }) => {
  const location = useLocation();

  const key = useMemo(() => {
    const searchParams = new URLSearchParams(location.search);
    return searchParams.get('key');
  }, [location])

  if (!key) {
    return <p>key is not specified</p>
  }

  return <button onClick={() => { context.backend.download(key) }}>Download</button>
}

export default DownloadPage;
