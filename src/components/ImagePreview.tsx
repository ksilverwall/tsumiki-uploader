import React from "react";
import "./ImagePreview.css";

const ImagePreview: React.FC<{ src: string; alt?: string; marked?: boolean }> = ({
  src,
  alt,
  marked,
}) => {
  const classes = ["image-preview"];
  if (marked) {
    classes.push("image-preview-marked");
  }

  return (
    <div className={classes.join(" ")}>
      <img src={src} alt={alt} />
    </div>
  );
};

export default ImagePreview;
