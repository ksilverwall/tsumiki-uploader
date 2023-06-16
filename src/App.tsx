import { useCallback, useEffect, useMemo, useReducer, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { uuidv7 } from "uuidv7";
import ImagePreview from "./ImagePreview";
import FileLoader from "./FileLoader";
import reducer from "./reducer";
import "./App.css";

function generateId<T extends string>(): T {
  const newId = uuidv7();
  return newId as T;
}

function App() {
  const navigate = useNavigate();
  const location = useLocation();

  const DEFAULT_GROUP_ID = useMemo(() => generateId<GroupId>(), []);

  const [viewGroupId, setViewGroupId] = useState<string>(DEFAULT_GROUP_ID);
  const [selectedItems, setSelectedItems] = useState<ItemId[]>([]);
  const [status, dispatch] = useReducer(reducer, {
    groups: {
      [DEFAULT_GROUP_ID]: {
        label: "無題",
        state: "EDITING",
        items: {},
      },
    },
  });

  useEffect(() => {
    const dict = Object.fromEntries(
      new URLSearchParams(location.search).entries()
    );
    setViewGroupId(dict["group"] ?? DEFAULT_GROUP_ID);
  }, [DEFAULT_GROUP_ID, location.search]);

  useEffect(() => {
    setSelectedItems([]);
  }, [viewGroupId]);

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

  const onLoad = useCallback(
    (groupId: GroupId, files: File[]) =>
      dispatch({
        type: "LOAD",
        groupId: groupId,
        items: files.map((f) => ({
          id: generateId<ItemId>(),
          file: f,
        })),
      }),
    []
  );

  const groupButtons = (
    <div style={{ display: "flex", flexDirection: "column" }}>
      {Object.keys(status.groups)
        .sort()
        .map((v) => (
          <button onClick={() => navigate(`?group=${v}`)}>{v}</button>
        ))}
      <button
        onClick={() =>
          dispatch({ type: "CREATE_GROUP", newGroupId: generateId() })
        }
      >
        +
      </button>
    </div>
  );

  const newGroupButton = viewGroupId ? (
    <button
      onClick={() => {
        dispatch({
          type: "CREATE_GROUP",
          newGroupId: generateId<GroupId>(),
          source: {
            groupId: viewGroupId,
            selected: selectedItems,
          },
        });
      }}
    >
      New Group
    </button>
  ) : null;

  const gallerySlots = viewGroupId
    ? {
        loaderPanel: (
          <FileLoader onLoaded={(file) => onLoad(viewGroupId, file)} />
        ),
        imagePanels: Object.keys(status.groups[viewGroupId].items).map(
          (id, idx) => (
            <div key={idx} onClick={() => onSelect(id)}>
              <ImagePreview
                file={status.groups[viewGroupId].items[id].file}
                marked={selectedItems.includes(id)}
              />
            </div>
          )
        ),
        galleryHeader: (
          <>
            <p>{status.groups[viewGroupId].label}</p>
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

  const archiveAllButton = <button>Archive All</button>;

  return (
    <>
      <header>
        <p>Tsumiki Uploader</p>
      </header>
      <div className="body-container">
        <aside className="sidebar">
          <div className="sidebar-contents">{groupButtons}</div>
          <div>{archiveAllButton}</div>
        </aside>
        <section className="main-section">{galleryView}</section>
      </div>
    </>
  );
}

export default App;
