import { useState } from "react";
import ImagePreview from "./ImagePreview";
import FileLoader from "./FileLoader";
import "./App.css";

function App() {
  const [selectedFiles, setSelectedFiles] = useState<File[]>([]);

  return (
    <div className="gallery-view">
      {selectedFiles.length > 0 ? (
        <>
          {selectedFiles.map((file, index) => (
            <ImagePreview key={index} file={file} />
          ))}
        </>
      ) : (
        <FileLoader onLoaded={setSelectedFiles} />
      )}
    </div>
  );
}

export default App;
