export function getData(){
    return [
        {
            "Path": "/flow.json",
            "IsDirectory": false
        },
        {
            "Path": "/flow1.json",
            "IsDirectory": false
        },
        {
            "Path": "/flow2.json",
            "IsDirectory": false
        },
        {
            "Path": "/flow3.json",
            "IsDirectory": false
        },
    ]
}

interface Dictionary {
    [key: string]: string;
}

export const getFiles: Dictionary = {
    "flow": `pub resource NFT {

        pub fun greet(): String {
      
          return "I'm NFT #"
      
            .concat(self.uuid.toString())
      
        }
      
      }
      
      
      pub fun main(): String {
      
        let nft <- create NFT()
      
        let greeting = nft.greet()
      
        destroy nft
      
        return greeting
      
      }`,
    "flow1" : `pub resource NFT {

        pub fun greet(): String {
      
          return "I'm flo1 #"
      
            .concat(self.uuid.toString())
      
        }
      
      }
      
      
      pub fun main(): String {
      
        let nft <- create NFT()
      
        let greeting = nft.greet()
      
        destroy nft
      
        return greeting
      
      }`,
    "flow2" : `pub resource NFT {

        pub fun greet(): String {
      
          return "I'm flow2 #"
      
            .concat(self.uuid.toString())
      
        }
      
      }
      
      
      pub fun main(): String {
      
        let nft <- create NFT()
      
        let greeting = nft.greet()
      
        destroy nft
      
        return greeting
      
      }`,
    "flow3" : `pub resource NFT {

        pub fun greet(): String {
      
          return "I'm flow3 #"
      
            .concat(self.uuid.toString())
      
        }
      
      }
      
      
      pub fun main(): String {
      
        let nft <- create NFT()
      
        let greeting = nft.greet()
      
        destroy nft
      
        return greeting
      
      }`

}