import { useCallback, useEffect, useState } from "react";
import GalleryViewLayout from "./GalleryViewLayout";
import FileLoader from "./FileLoader";
import ImagePreview from "./ImagePreview";
import { GenerateId } from "../app/libs";
import { ApplicationError } from "../app/repositories/ApplicationError";
import { ItemId, Item } from "../app/models";

export type GalleryViewProps = { label: string; items: { [key: ItemId]: Item } } & ({
  state: "EDITING";
  onCreateGroup: (ids: ItemId[]) => void;
  onUpdateItems: (items: { [key: ItemId]: Item }) => void;
} | {
  state: "ARCHIVING";
} | {
  state: "COMPLETE";
  url: URL;
} | {
  state: "FAILED";
  error: ApplicationError;
})

const GalleryView: React.FC<GalleryViewProps> = (props) => {
  const [selectedItems, setSelectedItems] = useState<ItemId[]>([]);
  const [items, setItems] = useState<{ [key: ItemId]: Item }>({});

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

  const onLoaded = useCallback((files: File[]) => {
    const newItems = {
      ...items,
      ...Object.fromEntries(
        files.map((f) => [
          GenerateId<ItemId>(),
          {
            file: f,
          },
        ])
      ),
    }
    setItems(newItems)
  }, [])

  useEffect(() => {
    if (props.state === "EDITING") {
      props.onUpdateItems(items)
    }
  }, [items]);

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
  ) : props.state === 'FAILED' ? (
    <div>
      {title}
      <p>{props.error.toString()}</p>
    </div>
  ) : (
    <div>
      {title}
      <p>{props.url.toString()}</p>
    </div>
  );

  const images = Object.keys(props.items).length > 0 ? (Object.keys(props.items).map((id, idx) => (
    <div key={idx} onClick={props.state === 'EDITING' ? () => onSelect(id) : undefined}>
      <ImagePreview
        src={props.items[id].file.type.startsWith("image/") ? URL.createObjectURL(props.items[id].file) : "xxxxx"}
        alt={props.items[id].file.name}
        marked={props.state === 'EDITING' ? selectedItems.includes(id) : false}
      />
    </div>
  ))) : [
    <p>Drag and drop files here</p>
  ]

  return (
    <FileLoader onLoaded={props.state === 'EDITING' ? onLoaded : undefined}>
      <GalleryViewLayout slots={{ header, images }} />
    </FileLoader>
  )
}

export default GalleryView;
