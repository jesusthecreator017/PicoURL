export interface ShortenResponse {
  short_url: string;
}

// Use just the shorten URL endpoint for now
export default async function shortenURL(url: string): Promise<ShortenResponse> {
  const res = await fetch("/api/shorten", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ url }),
  });

  const data = await res.json();

  if(!res.ok){
    throw new Error(data.error || 'Something went wrong');
  }

  return data;
}
