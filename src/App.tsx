import { useCallback, useEffect, useReducer, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import ImagePreview from "./ImagePreview";
import FileLoader from "./FileLoader";
import reducer from "./reducer";
import "./App.css";

function App() {
  const navigate = useNavigate();
  const location = useLocation();
  const DEFAULT_GROUP_ID = "0";

  const [groupView, setGroupView] = useState<string>(DEFAULT_GROUP_ID);
  const [status, dispatch] = useReducer(reducer, {
    nextArchiveId: 0,
    nextItemId: 0,
    groups: {
      [DEFAULT_GROUP_ID]: {items: []},
    },
  });

  useEffect(() => {
    const dict = Object.fromEntries(
      new URLSearchParams(location.search).entries()
    );
    setGroupView(dict["group"] ?? DEFAULT_GROUP_ID);
  }, [location.search]);

  const onSelect = useCallback(
    (index: number) => {
      dispatch({
        type: "MARK_ITEM",
        groupId: groupView,
        index: index,
        value: !status.groups[groupView].items[index].marked,
      });
    },
    [status, groupView]
  );

  const archiveButton = (
    <button
      onClick={() => {
        dispatch({ type: "ARCHIVE", groupId: groupView });
      }}
    >
      Archive
    </button>
  );

  const groupButtons = (
    <div style={{ display: "flex", flexDirection: "column" }}>
      {[Object.keys(status.groups)].sort().map((v) => (
        <button onClick={() => navigate(`?group=${v}`)}>{v}</button>
      ))}
    </div>
  );

  const renderItems = status.groups[groupView].items;

  const loaderPanel = (
    <FileLoader
      onLoaded={(files) =>
        dispatch({ type: "LOAD", groupId: groupView, files })
      }
    />
  );

  const imagePanels = renderItems.map((item, idx) => (
    <div key={idx} onClick={() => onSelect(item.id)}>
      <ImagePreview file={item.file} marked={item.marked} />
    </div>
  ));

  const galleryView = (
    <div className="gallery-view">
      <p>{groupView}</p>
      <div className="image-list">
        {status.groups[groupView].items.length > 0 ? (
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
