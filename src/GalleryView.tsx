import ImagePreview from "./ImagePreview";
import FileLoader from "./FileLoader";
import { useCallback, useEffect, useState } from "react";
import GalleryViewLayout from "./GalleryViewLayout";

const GalleryView: React.FC<{
  groupId: GroupId;
  group: Group;
  url?: string;
  onLoad: (id: GroupId, files: File[]) => void;
  onCreateGroup: (id: GroupId, selected: ItemId[]) => void;
}> = ({ groupId, group, url, onLoad, onCreateGroup }) => {
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
  }, [groupId]);

  if (!groupId || !group) {
    return null;
  }

  const slots = {
    header: (
      <>
        <p>{group.label}</p>
        {
          <button onClick={() => onCreateGroup(groupId, selectedItems)}>
            New Group
          </button>
        }
        {url ? <p>{url}</p> : null}
      </>
    ),
    images: [<FileLoader onLoaded={(file) => onLoad(groupId, file)} />].concat(Object.keys(group.items).map((id, idx) => (
      <div key={idx} onClick={() => onSelect(id)}>
        <ImagePreview
          file={group.items[id].file}
          marked={selectedItems.includes(id)}
        />
      </div>
    ))),
  }

  return <GalleryViewLayout slots={slots} />
};

export default GalleryView;
