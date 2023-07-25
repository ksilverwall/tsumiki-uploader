type Prop = { onLoaded?: (files: File[]) => void, className?: string }

const FileLoader: React.FC<React.PropsWithChildren<Prop>> = ({
  onLoaded,
  className,
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
    <div className={className} onDrop={handleDrop} onDragOver={handleDragOver}>
      {children}
    </div>
  );
};

export default FileLoader;
