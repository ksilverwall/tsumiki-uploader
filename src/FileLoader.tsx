type Prop = { onLoaded?: (files: File[]) => void }

const FileLoader: React.FC<React.PropsWithChildren<Prop>> = ({
  onLoaded,
  children,
}) => {
  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    const files = Array.from(e.dataTransfer.files);
    onLoaded && onLoaded(files);
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
  };

  return (
    <div onDrop={handleDrop} onDragOver={handleDragOver}>
      {children}
    </div>
  );
};

export default FileLoader;
