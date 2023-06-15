import React from "react";
import "./ImagePreview.css";

const ImagePreview: React.FC<{ file: File }> = ({ file }) => {
  if (file.type.startsWith("image/")) {
    return (
      <div className="image-preview">
        <img src={URL.createObjectURL(file)} alt={file.name} />
      </div>
    );
  } else {
    return <p>{file.name}</p>;
  }
};

export default ImagePreview;
