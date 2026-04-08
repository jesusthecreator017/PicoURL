import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import shortenURL from "./api";
import "./App.css";

function App() {
  const [url, setUrl] = useState<string>("");

  const { data, isPending, isError, error, mutate } = useMutation({
    mutationFn: shortenURL,
    onSuccess: () => setUrl(''),
  });
  const baseUrl = import.meta.env.DEV ? 'http://localhost:8080' : window.location.origin;
  const shortUrl = data ? `${baseUrl}/${data.short_url}` : null;

  return (
    <div>
      <h1>PicoURL</h1>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          mutate(url);
        }}
      >
        <input
          type="url"
          placeholder="https://example.com"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          required
        />
        <button type="submit" disabled={isPending}>
          {isPending ? "Shortening..." : "Shorten"}
        </button>
      </form>

      {shortUrl && (
        <p>
          Short URL:{" "}
          <a href={shortUrl} target="_blank" rel="noreferrer">
            {shortUrl}
          </a>
        </p>
      )}
      {isError && <p style={{ color: "red" }}>{error?.message ?? 'Something went wrong'}</p>}
    </div>
  );
}

export default App;
