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

  const onSelect = useCallback(
    (id: ItemId) => {
      dispatch({
        type: "MARK_ITEM",
        groupId: groupView,
        itemId: id,
        value: !status.groups[groupView].items[id].marked,
      });
    },
    [status, groupView]
  );

  const archiveButton = (
    <button
      onClick={() => {
        dispatch({
          type: "ARCHIVE",
          newGroupId: generateId<GroupId>(),
          groupId: groupView,
        });
      }}
    >
      Archive
    </button>
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
    (key, idx) => (
      <div key={idx} onClick={() => onSelect(key)}>
        <ImagePreview
          file={status.groups[groupView].items[key].file}
          marked={status.groups[groupView].items[key].marked}
        />
      </div>
    )
  );

  const galleryView = (
    <div className="gallery-view">
      <p>{groupView}</p>
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
          <div>{archiveButton}</div>
        </aside>
        <section className="main-section">{galleryView}</section>
      </div>
    </>
  );
}

export default App;
