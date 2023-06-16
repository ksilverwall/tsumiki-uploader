import { useCallback, useEffect, useReducer, useState } from "react";
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
  const DEFAULT_GROUP_ID = "0";

  const [groupView, setGroupView] = useState<string>(DEFAULT_GROUP_ID);
  const [selectedItems, setSelectedItems] = useState<ItemId[]>([]);
  const [status, dispatch] = useReducer(reducer, {
    groups: {
      [DEFAULT_GROUP_ID]: { items: {} },
    },
  });

  useEffect(() => {
    const dict = Object.fromEntries(
      new URLSearchParams(location.search).entries()
    );
    setGroupView(dict["group"] ?? DEFAULT_GROUP_ID);
  }, [location.search]);

  useEffect(() => {
    setSelectedItems([]);
  }, [groupView]);

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

  if (!status.groups[groupView]) {
    // FIXME: change url
    setGroupView(DEFAULT_GROUP_ID);
    return <p>Invalid GroupId</p>;
  }

  const newGroupButton = (
    <button
      onClick={() => {
        dispatch({
          type: "CREATE_GROUP",
          newGroupId: generateId<GroupId>(),
          source: {
            groupId: groupView,
            selected: selectedItems,
          },
        });
      }}
    >
      New Group
    </button>
  );

  const archiveAllButton = <button>Archive All</button>;

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

  const loaderPanel = (
    <FileLoader
      onLoaded={(files) =>
        dispatch({
          type: "LOAD",
          groupId: groupView,
          items: files.map((f) => ({
            id: generateId<ItemId>(),
            file: f,
          })),
        })
      }
    />
  );

  const imagePanels = Object.keys(status.groups[groupView].items).map(
    (id, idx) => (
      <div key={idx} onClick={() => onSelect(id)}>
        <ImagePreview
          file={status.groups[groupView].items[id].file}
          marked={selectedItems.includes(id)}
        />
      </div>
    )
  );

  const galleryView = (
    <div className="gallery-view">
      <div>
        <p>{groupView}</p>
        {newGroupButton}
      </div>
      <div className="image-list">
        {Object.keys(status.groups[groupView].items).length > 0 ? (
          <>
            <div>{loaderPanel}</div>
            {imagePanels.map((p) => (
              <div>{p}</div>
            ))}
          </>
        ) : (
          loaderPanel
        )}
      </div>
    </div>
  );

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
