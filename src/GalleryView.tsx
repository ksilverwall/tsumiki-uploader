import ImagePreview from "./ImagePreview";
import FileLoader from "./FileLoader";
import { useCallback, useEffect, useState } from "react";

const GalleryView: React.FC<{
  groupId: GroupId;
  group: Group;
  onLoad: (id: GroupId, files: File[]) => void;
  onCreateGroup: (id: GroupId, selected: ItemId[]) => void;
}> = ({ groupId, group, onLoad, onCreateGroup }) => {
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

  const newGroupButton = groupId ? (
    <button onClick={() => onCreateGroup(groupId, selectedItems)}>
      New Group
    </button>
  ) : null;

  const gallerySlots =
    groupId && group
      ? {
          loaderPanel: (
            <FileLoader onLoaded={(file) => onLoad(groupId, file)} />
          ),
          imagePanels: Object.keys(group.items).map((id, idx) => (
            <div key={idx} onClick={() => onSelect(id)}>
              <ImagePreview
                file={group.items[id].file}
                marked={selectedItems.includes(id)}
              />
            </div>
          )),
          galleryHeader: (
            <>
              <p>{group.label}</p>
              {newGroupButton}
            </>
          ),
        }
      : null;

  const galleryView = gallerySlots ? (
    <div className="gallery-view">
      <div>{gallerySlots.galleryHeader}</div>
      <div className="image-list">
        {gallerySlots.imagePanels.length > 0 ? (
          <>
            <div>{gallerySlots.loaderPanel}</div>
            {gallerySlots.imagePanels.map((p) => (
              <div>{p}</div>
            ))}
          </>
        ) : (
          gallerySlots.loaderPanel
        )}
      </div>
    </div>
  ) : null;

  return galleryView;
};

export default GalleryView;
