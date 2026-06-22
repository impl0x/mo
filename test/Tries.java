
import java.util.ArrayList;
import java.util.List;

public class Tries{
    public static void main(String[] args) {
        final List<String> setOfStrings=new ArrayList<>();
        setOfStrings.add("pqrs");
        setOfStrings.add("pprt");
        setOfStrings.add("psst");
        setOfStrings.add("qqrs");
        setOfStrings.add("pqrs");
    }
}

class Trie{
    final TrieNode root;
    public Trie(){
        this.root=new TrieNode();
    }
    public int query(String s){
        TrieNode current=root.next(s.charAt(0));
        for(int i=1;i<s.length();i++){
            current =current.next(s.charAt(i));
            if (current==null){
                return 0;
            }
        }
        return current.terminating;
    }
}

class TrieNode{
    int terminating;
    final TrieNode[] trieNodes=new TrieNode[26];

    public TrieNode next(final char c){
        return trieNodes[c-'a'];
    }


}
    