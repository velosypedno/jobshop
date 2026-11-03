export function getMachineColor(typeId: number): string {
  // Золотий кут (137.5°) забезпечує рівномірне розподілення кольорів по колірному колу
  const hue = (typeId * 137.5) % 360;
  // Насиченість 70%, світлість 45% гарантують яскравий, але читабельний колір на світлому фоні
  return `hsl(${hue}, 70%, 55%)`;
}

export function getMachineBackgroundColor(typeId: number): string {
  const hue = (typeId * 137.5) % 360;
  // Блідий варіант для фону (світлість 93%)
  return `hsl(${hue}, 70%, 93%)`;
}

export function buildMachineColorMap(machines: Array<{ type_id: number }> = []): Record<number, string> {
  const map: Record<number, string> = {};
  machines.forEach((m) => {
    map[m.type_id] = getMachineColor(m.type_id);
  });
  return map;
}
