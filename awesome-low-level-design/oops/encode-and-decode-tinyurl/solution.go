package encodedecodetinyurl

/*
LeetCode 535. Encode and Decode TinyURL

TinyURL is a URL shortening service where you enter a URL such as
https://leetcode.com/problems/design-tinyurl and it returns a short URL such as
https://tinyurl.com/4e9iAk. Design a class to encode a URL and decode a tiny URL.

There is no restriction on how your encode/decode algorithm should work. You just need
to ensure that a URL can be encoded to a tiny URL and the tiny URL can be decoded to the
original URL.

Implement the Codec class:
  - NewCodec() initializes the object.
  - Encode(longUrl string) string encodes a URL to a shortened URL.
  - Decode(shortUrl string) string decodes shortUrl back to the original URL.

Example 1:

	Input: url = "https://leetcode.com/problems/design-tinyurl"
	Output: "https://leetcode.com/problems/design-tinyurl"

	Explanation:
	Codec codec = new Codec();
	string shortUrl = codec.encode(url);  // returns a shortened, unique URL for `url`
	string longUrl = codec.decode(shortUrl); // returns "https://leetcode.com/problems/design-tinyurl"

Constraints:
  - 1 <= url.length <= 10^4
  - url is guaranteed to be a valid URL.

TODO:
 1. Add the fields you need on Codec below. The whole problem is really "design a
    reversible ID generator," which doesn't require any actual hashing or URL parsing:
      - a map from short code -> original long URL (for Decode), and
      - either a map from long URL -> short code (to make repeated Encode calls on the
        same URL idempotent, i.e. return the same short URL) OR a simple incrementing
        counter used as the short code, whichever you prefer.
 2. Implement NewCodec: initialize your map(s)/counter.
 3. Implement Encode(longUrl): if you've already encoded this exact URL before and want
    idempotent behavior, return the previously generated short URL. Otherwise, generate a
    new short code (e.g. a fixed prefix like "http://tinyurl.com/" + a base62 encoding of
    an incrementing counter, or just the counter itself as a string — the problem doesn't
    require a specific format), record the mapping, and return the short URL.
 4. Implement Decode(shortUrl): look up the code portion of shortUrl in your map and
    return the corresponding original long URL.
*/

type Codec struct {
	// TODO: add your fields here
}

func NewCodec() *Codec {
	// TODO: implement
	return &Codec{}
}

func (c *Codec) Encode(longUrl string) string {
	// TODO: implement
	return ""
}

func (c *Codec) Decode(shortUrl string) string {
	// TODO: implement
	return ""
}
