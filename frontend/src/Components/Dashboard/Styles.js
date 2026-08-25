import styled from "styled-components";

export const PageHeader = styled.div`
  margin-bottom: 28px;
`;

export const Title = styled.h1`
  font-family: var(--font-heading, 'Plus Jakarta Sans', sans-serif);
  font-size: 28px;
  font-weight: 700;
  color: #f8fafc;
  margin: 0;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #ffffff 0%, #9ca3af 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
`;

export const StatsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
`;

export const StatCard = styled.div`
  background: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  padding: 22px 24px;

  display: flex;
  align-items: center;
  gap: 18px;

  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, #38bdf8, #818cf8);
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  &:hover {
    border-color: rgba(56, 189, 248, 0.3);
    box-shadow: 0 12px 30px rgba(0, 0, 0, 0.4), 0 0 20px rgba(56, 189, 248, 0.15);
    transform: translateY(-4px);
    background: rgba(15, 23, 42, 0.85);

    &::before {
      opacity: 1;
    }
  }
`;

export const StatIcon = styled.div`
  width: 52px;
  height: 52px;

  display: flex;
  align-items: center;
  justify-content: center;

  background: linear-gradient(
    135deg,
    rgba(56, 189, 248, 0.15),
    rgba(129, 140, 248, 0.15)
  );
  border: 1px solid rgba(56, 189, 248, 0.2);

  border-radius: 14px;
  transition: transform 0.3s ease;

  ${StatCard}:hover & {
    transform: scale(1.1) rotate(4deg);
  }
`;

export const StatInfo = styled.div`
  flex: 1;
`;

export const StatValue = styled.div`
  font-family: var(--font-heading, 'Plus Jakarta Sans', sans-serif);
  font-size: 28px;
  font-weight: 800;
  color: #f8fafc;
  margin-bottom: 2px;
  letter-spacing: -0.02em;
`;

export const StatLabel = styled.div`
  font-family: var(--font-body, 'Inter', sans-serif);
  font-size: 13.5px;
  color: #94a3b8;
  font-weight: 500;
`;

export const Section = styled.div`
  background: rgba(15, 23, 42, 0.65);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  padding: 28px;

  @media (max-width: 768px) {
    padding: 18px;
  }
`;

export const SectionHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 16px;
`;

export const SectionTitle = styled.h2`
  font-family: var(--font-heading, 'Plus Jakarta Sans', sans-serif);
  font-size: 18px;
  font-weight: 700;
  color: #f8fafc;
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0;
`;

export const ThresholdBadge = styled.div`
  padding: 6px 14px;
  background: rgba(244, 63, 94, 0.12);
  border: 1px solid rgba(244, 63, 94, 0.25);
  border-radius: 20px;
  color: #fb7185;
  font-size: 12.5px;
  font-weight: 600;
  letter-spacing: 0.02em;
`;

export const TableWrapper = styled.div`
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  overflow: hidden;
  margin-top: 20px;
`;

export const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
`;

export const Th = styled.th`
  padding: 14px 18px;
  background: rgba(30, 41, 59, 0.5);
  color: #94a3b8;
  font-size: 13px;
  font-weight: 600;
  text-align: left;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  text-transform: uppercase;
  letter-spacing: 0.05em;
`;

export const Td = styled.td`
  padding: 16px 18px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  color: #f1f5f9;
  font-size: 14px;
  font-weight: 500;
`;

export const Icon = styled.img`
  width: 20px;
  height: 20px;
  filter: invert(80%) sepia(20%) saturate(1000%) hue-rotate(170deg);
`;
