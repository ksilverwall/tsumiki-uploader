import "./FileLoader.css"

const FileLoader: React.FC<{ onLoaded: (files: File[]) => void }> = ({
  onLoaded: onSelected,
}) => {
  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    const files = Array.from(e.dataTransfer.files);
    onSelected && onSelected(files);
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
  };

  return (
    <div
      className="file-selector"
      onDrop={handleDrop}
      onDragOver={handleDragOver}
    >
      <p>Drag and drop files here</p>
    </div>
  );
};

export default FileLoader;
