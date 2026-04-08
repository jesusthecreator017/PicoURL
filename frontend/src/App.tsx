import { useState, useEffect } from "react";
import { useMutation } from "@tanstack/react-query";
import shortenURL from "./api";
import "./App.css";

function Typewriter({ text, onDone }: { text: string; onDone?: () => void }) {
  const [displayed, setDisplayed] = useState("");

  useEffect(() => {
    setDisplayed("");
    let i = 0;
    const interval = setInterval(() => {
      i++;
      setDisplayed(text.slice(0, i));
      if (i >= text.length) {
        clearInterval(interval);
        onDone?.();
      }
    }, 35);
    return () => clearInterval(interval);
  }, [text, onDone]);

  return (
    <span>
      {displayed}
      <span className="typewriter-cursor">|</span>
    </span>
  );
}

function App() {
  const [url, setUrl] = useState("");
  const [copied, setCopied] = useState(false);
  const [typing, setTyping] = useState(false);

  const { data, isPending, isError, error, mutate } = useMutation({
    mutationFn: shortenURL,
    onSuccess: () => {
      setUrl("");
      setCopied(false);
      setTyping(true);
    },
  });

  const baseUrl = import.meta.env.DEV
    ? "http://localhost:8080"
    : window.location.origin;
  const shortUrl = data ? `${baseUrl}/${data.short_url}` : null;

  const handleCopy = async () => {
    if (!shortUrl) return;
    await navigator.clipboard.writeText(shortUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <>
      <h1 className="title">
        Pico<span>URL</span>
      </h1>

      <div className="cassette-wrapper">
        <div className={`cassette${isPending ? " spinning" : ""}`}>
          {/* Corner screws */}
          <div className="screw screw-tl" />
          <div className="screw screw-tr" />
          <div className="screw screw-bl" />
          <div className="screw screw-br" />

          {/* Label with retro stripes */}
          <div className="label">
            <div className="label-stripes">
              <div className="stripe stripe-cream" />
              <div className="stripe stripe-sandy" />
              <div className="stripe stripe-peach" />
              <div className="stripe stripe-salmon" />
              <div className="stripe stripe-verdigris" />
            </div>

            {/* Input overlaid on cream area */}
            <div className="input-area">
              <div className="input-label">paste your url</div>
              <form
                className="input-row"
                onSubmit={(e) => {
                  e.preventDefault();
                  mutate(url);
                }}
              >
                <input
                  className="url-input"
                  type="url"
                  placeholder="https://example.com"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  required
                />
                <button
                  className="shorten-btn"
                  type="submit"
                  disabled={isPending}
                >
                  {isPending ? "..." : "REC"}
                </button>
              </form>
            </div>

            {/* Tape window overlaid on stripes */}
            <div className="tape-window">
              <div className="reel" />
              <div className="tape-center" />
              <div className="reel" />
            </div>
          </div>

          {/* Bottom section — tape sticker always visible */}
          <div className="cassette-bottom">
            <div className="tape-sticker">
              {shortUrl ? (
                <>
                  <a href={shortUrl} target="_blank" rel="noreferrer">
                    {typing ? (
                      <Typewriter
                        text={shortUrl}
                        onDone={() => setTyping(false)}
                      />
                    ) : (
                      shortUrl
                    )}
                  </a>
                  {!typing && (
                    <button
                      className={`copy-btn${copied ? " copied" : ""}`}
                      onClick={handleCopy}
                    >
                      {copied ? "copied" : "copy"}
                    </button>
                  )}
                </>
              ) : (
                <span className="tape-sticker-placeholder">
                  {isPending ? "recording..." : "side a"}
                </span>
              )}
            </div>

            {isError && (
              <p className="error-msg">
                {error?.message ?? "Something went wrong"}
              </p>
            )}

            <div className="bottom-dots-row">
              <div className="bottom-dots">
                <div className="bottom-dot dot-orange" />
                <div className="bottom-dot dot-orange" />
              </div>
              <div className="bottom-dots">
                <div className="bottom-dot dot-orange" />
                <div className="bottom-dot dot-grey" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

export default App;
