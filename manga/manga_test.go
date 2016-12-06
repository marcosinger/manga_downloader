package manga

import "testing"

type Node struct {
	attributes map[string]string
}

func (n Node) Attr(attr string) (string, bool) {
	v, exists := n.attributes[attr]
	return v, exists
}

func TestBuildImageName(t *testing.T) {
	var tests = []struct {
		node Node
		name string
		src  string
	}{
		{
			Node{attributes: map[string]string{"alt": "Onepunch-Man 91 - Page 1", "src": "http://cdn.com/onepunch-man.jpg"}},
			"Onepunch-Man-Chap-00091-Pg-00001.jpg",
			"http://cdn.com/onepunch-man.jpg",
		},
		{
			Node{attributes: map[string]string{"alt": "Hajime no Ippo 3 - Page 20", "src": "http://cdn.com/hajime-no-ippo.jpg"}},
			"Hajime-no-Ippo-Chap-00003-Pg-00020.jpg",
			"http://cdn.com/hajime-no-ippo.jpg",
		},
		{
			Node{attributes: map[string]string{"alt": "Fairy Tail 512 - Page 21", "src": "http://cdn.com/fairy-tail.jpg"}},
			"Fairy-Tail-Chap-00512-Pg-00021.jpg",
			"http://cdn.com/fairy-tail.jpg",
		},
		{
			Node{attributes: map[string]string{"alt": "Gamon The Demolition Man 650 - Page 810", "src": "http://cdn.com/gamon-the-demolition-man.jpg"}},
			"Gamon-The-Demolition-Man-Chap-00650-Pg-00810.jpg",
			"http://cdn.com/gamon-the-demolition-man.jpg",
		},
	}

	for _, test := range tests {
		name, src := buildImageName(test.node)

		if name != test.name {
			t.Errorf("expected %s got %s", test.name, name)
		}

		if src != test.src {
			t.Errorf("expected %s got %s", test.src, src)
		}
	}
}
