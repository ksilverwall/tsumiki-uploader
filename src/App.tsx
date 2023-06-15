import { useCallback, useReducer } from "react";
import ImagePreview from "./ImagePreview";
import FileLoader from "./FileLoader";
import "./App.css";
import reducer from "./reducer";

function App() {
  const [status, dispatch] = useReducer(reducer, {
    nextArchiveId: 0,
    nextItemId: 0,
    items: [],
  });

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
          <button>{v}</button>
        ))}
    </div>
  );

  const galleryView = (
    <div className="gallery-view">
      {status.items.length > 0 ? (
        <>
          {status.items
            .filter((item) => item.archiveiId === undefined)
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
