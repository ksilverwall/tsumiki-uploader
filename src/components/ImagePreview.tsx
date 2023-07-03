import React from "react";
import "./ImagePreview.css";

const ImagePreview: React.FC<{ file: File; marked: boolean }> = ({
  file,
  marked,
}) => {
  if (file.type.startsWith("image/")) {
    const classes = ["image-preview"];
    if (marked) {
      classes.push("image-preview-marked");
    }

    return (
      <div className={classes.join(" ")}>
        <img src={URL.createObjectURL(file)} alt={file.name} />
      </div>
    );
  } else {
    return <p>{file.name}</p>;
  }
};

export default ImagePreview;
