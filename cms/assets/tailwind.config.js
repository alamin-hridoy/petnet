module.exports = {
  purge: [],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      fontFamily: {
      sans:[
        '"Proxima Nova"'
      ]
    },
    colors: {
        petnetblue:'#1A2791',
        petnetyellow:'#FDD100',
        petnetgray:'#F6F7FF',
        petnetpink:'#D12C7F',
        petnetpurple:'#6633CC',
        petnetgreen:'#46B746',
        petnetorange:'#F76F34',
        petnetlightblue:'#05ACE5',
        petnetheader:'#dfe0ea',
        petnetheadertext:"#BCC3C8",
        petnetlightyellow: "#FEF1B3",
        petnetdarkblue:"#041250",
        petnetred: "#FF0E0E",
        petnetlighterblue:"#B4E6F7",
        petnettextgray:"#3D3D3D",
        petnetstatuspink:"#E895BF",
        petnetinput:"#F3F4FA",
        petnetslightgray:"#DADADA",
        petnetlightgray:"#353535"
      },
    },
  },
  variants: {
    extend: {
      backgroundColor: ['checked'],
      borderColor: ['checked'],
    },
  },
  plugins: [],
}
