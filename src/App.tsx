import { useCallback, useEffect, useReducer, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import ImagePreview from "./ImagePreview";
import FileLoader from "./FileLoader";
import reducer from "./reducer";
import "./App.css";

function App() {
  const navigate = useNavigate();
  const location = useLocation();

  const [groupView, setGroupView] = useState<number>();
  const [status, dispatch] = useReducer(reducer, {
    nextArchiveId: 0,
    nextItemId: 0,
    items: [],
  });

  useEffect(() => {
    const dict = Object.fromEntries(
      new URLSearchParams(location.search).entries()
    );
    setGroupView(dict["group"] ? parseInt(dict["group"]) : undefined);
  }, [location.search]);

  const onSelect = useCallback(
    (index: number) => {
      dispatch({
        type: "MARK_ITEM",
        index: index,
        value: !status.items[index].marked,
      });
    },
    [status]
  );

  const archiveButton = (
    <button
      onClick={() => {
        dispatch({ type: "ARCHIVE" });
      }}
    >
      Archive
    </button>
  );

  const groupButtons = (
    <div style={{ display: "flex", flexDirection: "column" }}>
      {[
        ...new Set(
          status.items
            .map((item) => item.archiveiId)
            .filter((v) => v !== undefined)
        ),
      ]
        .sort()
        .map((v) => (
          <button onClick={() => navigate(`?group=${v}`)}>{v}</button>
        ))}
    </div>
  );

  const galleryView = (
    <div className="gallery-view">
      <p>{groupView}</p>
      <div className="image-list">
        {status.items.length > 0 ? (
          <>
            {status.items
              .filter((item) => item.archiveiId === groupView)
              .map((item, idx) => (
                <div key={idx} onClick={() => onSelect(item.id)}>
                  <ImagePreview file={item.file} marked={item.marked} />
                </div>
              ))}
          </>
        ) : (
          <FileLoader onLoaded={(files) => dispatch({ type: "LOAD", files })} />
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
