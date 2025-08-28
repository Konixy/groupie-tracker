// Script de test pour la logique de recherche
// À exécuter dans la console du navigateur

// Données de test (simulation des artistes)
const testArtists = [
  {
    id: 1,
    name: "Queen",
    members: ["Freddie Mercury", "Brian May", "Roger Taylor", "John Deacon"],
    creationDate: 1970,
    firstAlbum: "1973",
  },
  {
    id: 2,
    name: "The Beatles",
    members: [
      "John Lennon",
      "Paul McCartney",
      "George Harrison",
      "Ringo Starr",
    ],
    creationDate: 1960,
    firstAlbum: "1963",
  },
  {
    id: 3,
    name: "Pink Floyd",
    members: ["Roger Waters", "David Gilmour", "Nick Mason", "Richard Wright"],
    creationDate: 1965,
    firstAlbum: "1967",
  },
];

// Test de la logique de recherche
function testSearchLogic() {
  console.log("🧪 Test de la logique de recherche");
  console.log("=====================================");

  // Test 1: Recherche par nom d'artiste
  console.log("\n1️⃣ Test recherche par nom d'artiste:");
  const artistResults = testArtists.filter((artist) =>
    artist.name.toLowerCase().includes("queen")
  );
  console.log(
    "Recherche 'queen':",
    artistResults.map((a) => a.name)
  );

  // Test 2: Recherche par membre
  console.log("\n2️⃣ Test recherche par membre:");
  const memberResults = [];
  testArtists.forEach((artist) => {
    artist.members.forEach((member) => {
      if (member.toLowerCase().includes("freddie")) {
        memberResults.push({
          member: member,
          artist: artist.name,
        });
      }
    });
  });
  console.log("Recherche 'freddie':", memberResults);

  // Test 3: Recherche par date de création
  console.log("\n3️⃣ Test recherche par date de création:");
  const creationResults = testArtists.filter((artist) =>
    artist.creationDate.toString().includes("1970")
  );
  console.log(
    "Recherche '1970':",
    creationResults.map((a) => `${a.name} (${a.creationDate})`)
  );

  // Test 4: Recherche par premier album
  console.log("\n4️⃣ Test recherche par premier album:");
  const albumResults = testArtists.filter((artist) =>
    artist.firstAlbum.includes("1963")
  );
  console.log(
    "Recherche '1963':",
    albumResults.map((a) => `${a.name} (${a.firstAlbum})`)
  );

  console.log("\n✅ Tests terminés !");
}

// Exécuter les tests
testSearchLogic();

// Fonction utilitaire pour tester la recherche en temps réel
function simulateSearch(query) {
  console.log(`\n🔍 Simulation de recherche pour: "${query}"`);

  const results = [];

  // Recherche par nom d'artiste
  testArtists.forEach((artist) => {
    if (artist.name.toLowerCase().includes(query.toLowerCase())) {
      results.push({
        type: "artist",
        artist: artist.name,
        score: 1000,
      });
    }
  });

  // Recherche par membre
  testArtists.forEach((artist) => {
    artist.members.forEach((member) => {
      if (member.toLowerCase().includes(query.toLowerCase())) {
        results.push({
          type: "member",
          artist: artist.name,
          member: member,
          score: 900,
        });
      }
    });
  });

  // Recherche par date
  testArtists.forEach((artist) => {
    if (artist.creationDate.toString().includes(query)) {
      results.push({
        type: "creation_date",
        artist: artist.name,
        year: artist.creationDate,
        score: 700,
      });
    }
  });

  console.log("Résultats:", results);
  return results;
}

// Exemples d'utilisation
console.log("\n🎯 Exemples de recherche:");
simulateSearch("queen");
simulateSearch("freddie");
simulateSearch("1970");
simulateSearch("1963");
