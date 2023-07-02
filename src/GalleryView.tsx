import { useCallback, useEffect, useState } from "react";
import GalleryViewLayout from "./GalleryViewLayout";
import FileLoader from "./FileLoader";
import ImagePreview from "./ImagePreview";

export type GalleryViewProps = Group & ({
  state: "EDITING";
  onCreateGroup: (ids: ItemId[]) => void;
  onLoad: (files: File[]) => void;
} | {
  state: "ARCHIVING";
} | {
  state: "COMPLETE";
  url: string;
})

const GalleryView: React.FC<GalleryViewProps> = (props) => {
  const [selectedItems, setSelectedItems] = useState<ItemId[]>([]);

  const onSelect = useCallback(
    (id: ItemId) => {
      if (selectedItems.includes(id)) {
        setSelectedItems(selectedItems.filter((v) => v !== id));
      } else {
        setSelectedItems(selectedItems.concat(id));
      }
    },
    [selectedItems]
  );

  useEffect(() => {
    setSelectedItems([]);
  }, [props]);

  const title = <p>{props.label}</p>

  const header = props.state === 'EDITING' ? (
    <div>
      {title}
      <button onClick={() => props.onCreateGroup(selectedItems)}>New Group</button>
    </div>
  ) : props.state === 'ARCHIVING' ? (
    <div>
      {title}
    </div>
  ) : (
    <div>
      {title}
      <p>{props.url}</p>
    </div>
  );

  const images = Object.keys(props.items).length > 0 ? (Object.keys(props.items).map((id, idx) => (
    <div key={idx} onClick={props.state === 'EDITING' ? () => onSelect(id) : undefined}>
      <ImagePreview
        file={props.items[id].file}
        marked={props.state === 'EDITING' ? selectedItems.includes(id) : false}
      />
    </div>
  ))) : [
    <p>Drag and drop files here</p>
  ]

  return (
    <FileLoader onLoaded={props.state === 'EDITING' ? props.onLoad : undefined}>
      <GalleryViewLayout slots={{ header, images }} />
    </FileLoader>
  )
}

export default GalleryView;
